package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/sillsdev/go-aeneas/audiogenerators"
	"github.com/sillsdev/go-aeneas/datatypes"
)

var (
	logLevel = 0
	batch    = ""
)

func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	results := make([]string, 12)

	err = nil
	for err == nil {
		str, err := reader.ReadString('\n')
		if err != io.EOF && err != nil {
			return nil, err
		}
		if str != "" {
			results = append(results, str)
		}
	}

	return results, nil
}

type PhraseReadResults struct {
	phrases []*datatypes.Phrase
	err     error
}

func readPhrasesFromFile(filename string, phraseResults chan PhraseReadResults) {
	phrases, err := readFileLines(filename)
	if err != nil {
		phraseResults <- PhraseReadResults{nil, err}
		return
	}
	parsedPhrases := make([]*datatypes.Phrase, len(phrases))
	for _, phrase := range phrases {
		parsedPhrase, err := datatypes.ParsePhrase(phrase)
		if err != nil {
			phraseResults <- PhraseReadResults{nil, err}
			return
		}
		parsedPhrases = append(parsedPhrases, parsedPhrase)
	}

	phraseResults <- PhraseReadResults{parsedPhrases, nil}
}

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

	phrases := make(chan PhraseReadResults)
	go readPhrasesFromFile(tpv.Task.PhraseFilename, phrases)

	tpv.Println("Wave Filepath:", <-wavs)
	phraseResults := <-phrases
	if phraseResults.err != nil {
		tpv.Println("Error parsing phrases:", phraseResults.err)
	}
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

	results := make(chan string)
	var generator datatypes.AudioGenerator = audiogenerators.GetEspeakGenerator()

	tempDir := createTempDir()

	for _, task := range tasks {
		go processTask(results, task, &generator, tempDir)
	}

	for range tasks {
		fmt.Println(<-results)
	}

}
