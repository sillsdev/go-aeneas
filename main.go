package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/sillsdev/go-aeneas/audiogenerators"
	"github.com/sillsdev/go-aeneas/datatypes"
	"github.com/sillsdev/go-aeneas/dtw"
	"github.com/sillsdev/go-aeneas/mfcc"
)

var (
	logLevel       = 0
	batch          = ""
	plot           = false
	listGenerators = false
	generator      = ""
)

func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewScanner(file)

	results := make([]string, 0)

	for reader.Scan() {
		if text := reader.Text(); text != "" {
			results = append(results, text)
		}
	}

	if err := reader.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

type PhraseReadResults struct {
	phrase *datatypes.Phrase
	err    error
}

/**
 * Reads phrases from file, returning a channel with parsed phrases
 *
 * Each individual line may fail to parse if the file is malformed,
 * and this error is passed along per line in the channel returned
 *
 * Closes the channel provided as input
 */
func readPhrasesFromFile(filename string, phraseResults chan<- PhraseReadResults) {
	defer close(phraseResults)

	phrases, err := readFileLines(filename)
	if err != nil {
		phraseResults <- PhraseReadResults{nil, err}
		return
	}

	for _, phrase := range phrases {
		parsedPhrase, err := datatypes.ParsePhrase(phrase)
		if err != nil {
			phraseResults <- PhraseReadResults{nil, err}
			return
		} else {
			phraseResults <- PhraseReadResults{parsedPhrase, nil}
		}
	}
}

type PhraseWav struct {
	phrase            *datatypes.Phrase
	phraseWavFilePath string
}

type PhraseWavResults struct {
	phraseWav *PhraseWav
	err       error
}

/**
 * Management function for spawning go routines to generate WAV files and manage the channel to do so
 *
 * Closes the channel provided as input
 */
func generateWavFilesForPhrases(tpv *datatypes.TaskProcessVariables, phraseOrder chan<- *datatypes.Phrase, phraseResults <-chan PhraseReadResults, phrasesGenerated chan<- PhraseWavResults) {
	defer close(phrasesGenerated)
	defer close(phraseOrder)

	var wg sync.WaitGroup

	for phraseResult := range phraseResults {
		if phraseResult.err != nil {
			phrasesGenerated <- PhraseWavResults{nil, phraseResult.err}
			return
		}

		phraseOrder <- phraseResult.phrase
		phrase := phraseResult.phrase
		//fmt.Println("Phrase Recieved: ", phraseOrder) // Read the value from the channel

		wg.Add(1)
		go func() {
			defer wg.Done()
			generateWavFileForPhrase(tpv, phrase, phrasesGenerated)
		}()
	}

	wg.Wait()
}

func generateWavFileForPhrase(tpv *datatypes.TaskProcessVariables, phrase *datatypes.Phrase, phrasesGenerated chan<- PhraseWavResults) {
	err := (*tpv.Generator).GenerateAudioFile(tpv.Parameters, phrase, tpv.GetPhraseFilePath(phrase.PhraseIndex))

	if err != nil {
		phrasesGenerated <- PhraseWavResults{nil, err}
		return
	}

	phrasesGenerated <- PhraseWavResults{&PhraseWav{phrase, tpv.GetPhraseFilePath(phrase.PhraseIndex)}, nil}
}

type GeneratedMfcCoefficients struct {
	phraseAndWav PhraseWav
	mfccResult   *[][]float64
}

type MfccResults struct {
	mfcc *GeneratedMfcCoefficients
	err  error
}

func generateMfccForWavFiles(tpv *datatypes.TaskProcessVariables, phrasesGenerated <-chan PhraseWavResults, mfccResults chan<- MfccResults) {
	defer close(mfccResults)

	var wg sync.WaitGroup

	for phraseResult := range phrasesGenerated {
		if phraseResult.err != nil {
			mfccResults <- MfccResults{nil, phraseResult.err}
			return
		}

		phraseAndWav := phraseResult.phraseWav
		wg.Add(1)
		go func() {
			defer wg.Done()
			// do your mfcc, then write to mfccResults
			results, err := mfcc.GenerateMfcc(phraseAndWav.phraseWavFilePath)
			if err != nil {
				mfccResults <- MfccResults{nil, err}
				return
			}
			mfccResults <- MfccResults{&GeneratedMfcCoefficients{*phraseAndWav, &results}, nil}
		}()
	}
	wg.Wait()
}

/**
 * Process task pipeline
 *
 * Goes through initializing and managing each step:
 * - Convert input audio to WAV
 * - Parse provided phrases
 * - Generate WAV files from parsed phrases
 * - Prepare for MFCC
 * - Prepare for DTW
 */
func processTask(results chan string, task *datatypes.Task, generator *datatypes.AudioGenerator, tempDir string) {
	tpv := datatypes.NewTaskProcessVariables(task, generator, tempDir)
	defer func() {
		results <- tpv.GetFinalLogs()
	}()

	if len(task.Description) > 0 {
		tpv.Println("")
		tpv.Println("*** ", tpv.Task.Description, " ***")
		tpv.Println("")
	}
	tpv.Println("Audio   : ", tpv.Task.AudioFilename)
	tpv.Println("Phrase  : ", tpv.Task.PhraseFilename)
	tpv.Println("Output  : ", tpv.Task.OutputFilename)
	tpv.Println("Parameters : ", tpv.Parameters)

	wavs := make(chan string)
	go convertWav(wavs, tpv)

	phraseReads := make(chan PhraseReadResults)
	go readPhrasesFromFile(tpv.Task.PhraseFilename, phraseReads)
	phrasesWithFiles := make(chan PhraseWavResults)
	phraseOrder := make(chan *datatypes.Phrase)

	go generateWavFilesForPhrases(tpv, phraseOrder, phraseReads, phrasesWithFiles)
	mfccPhraseResults := make(chan MfccResults)
	go generateMfccForWavFiles(tpv, phrasesWithFiles, mfccPhraseResults)

	//fmt.Println("Number of Ordered Phrases Processed: ", len(phraseOrder))

	mfccPhrasesMap := make(map[string]*MfccResults)

	timeOffsetFloat, err := strconv.ParseFloat(tpv.Parameters.Get("is_audio_file_detect_head_max"), 64)
	if err != nil {
		tpv.Println("Error: ", err)
		return
	}
	tpv.Println("Initial Time Offset: ", timeOffsetFloat)
	timeOffset := (int)(timeOffsetFloat * 22050)

	file, err := os.Create(tpv.Task.OutputFilename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	splitDescription := strings.Split(tpv.Task.Description, " ")
	book := strings.TrimSpace(splitDescription[0])
	chapter := strings.TrimSpace(splitDescription[1])

	content := fmt.Sprintf(`\id %s
\c %s
\level phrase
\separators . ? ! : ; ,
`, book, chapter)

	// Write content to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	tpv.Println("Timing File created and written successfully.")

	mfccResultsChan := make(chan error)
	go func() {
		mfccInputResults, err := mfcc.GenerateMfcc(<-wavs)
		if err != nil {
			mfccResultsChan <- err
		} else {
			tpv.MfccInputResults = mfccInputResults
			mfccResultsChan <- nil
		}
	}()

	if err := <-mfccResultsChan; err != nil {
		tpv.Println("Error handling MFCC ", err)
	}

	for phrase := range phraseOrder {
		nextPhrase := phrase.PhraseIndex

		ok := false
		_, ok = mfccPhrasesMap[nextPhrase]

		for !ok {
			val := <-mfccPhraseResults
			mfccPhrasesMap[val.mfcc.phraseAndWav.phrase.PhraseIndex] = &val
			_, ok = mfccPhrasesMap[nextPhrase]
		}

		tpv.Println("Handling phrase MFCC/DTW: ", mfccPhrasesMap[nextPhrase].mfcc.phraseAndWav.phrase.PhraseIndex)
		oldTimeOffset := timeOffset
		timeOffset = dtw.RunDtw(tpv.MfccInputResults, *mfccPhrasesMap[nextPhrase].mfcc.mfccResult, timeOffset)

		temp := fmt.Sprintf("%d\t%d\t%s\n", oldTimeOffset/22050, timeOffset/22050, mfccPhrasesMap[nextPhrase].mfcc.phraseAndWav.phrase.PhraseIndex)
		file.WriteString(temp)

		mfccPhrasesMap[nextPhrase] = nil
		// update time offset
	}

	tpv.Println("Done with ", tpv.Task.Description, "!")
}

func createTempDir() string {
	TempDir, err := os.MkdirTemp("", "goaeneas")
	if err != nil {
		log.Fatal(err)
	}
	return TempDir
}

func convertWav(wavs chan<- string, tpv *datatypes.TaskProcessVariables) {
	filepath := tpv.GetWavFilepath()
	out, _ := exec.Command("ffmpeg", "-i", tpv.Task.AudioFilename, "-acodec", "pcm_s16le", "-ac", "1", "-ar", "22050", filepath).CombinedOutput() //Swapped to a sample rate of 22050 from 16000
	wavs <- filepath
	tpv.Println("ffmpeg output : ", string(out))
}

func main() {
	processArguments()

	audioGens := audiogenerators.GetAudioGenerators()

	if listGenerators {
		fmt.Println("Audio generators available:")
		for _, generator := range audioGens {
			fmt.Printf("\t%s\n", generator.GetName())
		}

		os.Exit(0)
	}

	tasks := []*datatypes.Task{}
	if len(batch) > 0 {
		//fmt.Println("Batch file:", batch)
		content, err := os.ReadFile(batch)
		if err != nil {
			log.Fatal("Error while reading batch file", err)
		}

		err = json.Unmarshal(content, &tasks)
		if err != nil {
			log.Fatal("Error parsing batch json file", err)
		}
	} else if len(os.Args) >= 5 {
		task := &datatypes.Task{
			Description:    "",
			AudioFilename:  os.Args[1],
			PhraseFilename: os.Args[2],
			Parameters:     os.Args[3],
			OutputFilename: os.Args[4],
		}
		tasks = append(tasks, task)
	}

	var finalAudioGenerator *datatypes.AudioGenerator = nil
	for _, availableGen := range audioGens {
		if availableGen.GetName() == generator {
			finalAudioGenerator = &availableGen
		}
	}

	fmt.Printf("Using audio generator %s\n", (*finalAudioGenerator).GetName())

	results := make(chan string)
	tempDir := createTempDir()

	for _, task := range tasks {
		go processTask(results, task, finalAudioGenerator, tempDir)
	}

	for range tasks {
		fmt.Println(<-results)
	}

}
