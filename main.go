package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/sillsdev/go-aeneas/audiogenerators"
	"github.com/sillsdev/go-aeneas/datatypes"
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
func generateWavFilesForPhrases(tpv *datatypes.TaskProcessVariables, phraseResults <-chan PhraseReadResults, phrasesGenerated chan<- PhraseWavResults) {
	defer close(phrasesGenerated)

	var wg sync.WaitGroup

	for phraseResult := range phraseResults {
		if phraseResult.err != nil {
			phrasesGenerated <- PhraseWavResults{nil, phraseResult.err}
			return
		}

		phrase := phraseResult.phrase

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
	go generateWavFilesForPhrases(tpv, phraseReads, phrasesWithFiles)

	tpv.Println("Logs for generated phrases:")
	for phraseLogItem := range phrasesWithFiles {
		if phraseLogItem.err != nil {
			tpv.Println("\tError generating phrase: ", phraseLogItem.err)
		} else {
			tpv.Println("\tFile successfully generated! ", phraseLogItem.phraseWav.phraseWavFilePath)
		}
	}

	mfccResultsChan := make(chan error)
	go func() {
		mfccResults, err := mfcc.GenerateMfcc(<-wavs)
		if err != nil {
			mfccResultsChan <- err
		} else {
			tpv.MfccResults = mfccResults
			mfccResultsChan <- nil
		}
	}()

	//
	if err := <-mfccResultsChan; err != nil {
		//
		tpv.Println("Error handling MFCC ", err)
	}

	if plot {
		mfcc.PlotMFCC(tpv.MfccResults)
	}
	tpv.Println("Done!")
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
	out, _ := exec.Command("ffmpeg", "-i", tpv.Task.AudioFilename, "-acodec", "pcm_s16le", "-ac", "1", "-ar", "16000", filepath).CombinedOutput()
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
		fmt.Println("Batch file:", batch)
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
