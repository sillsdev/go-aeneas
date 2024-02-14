package main

import (
	"encoding/json"
	"fmt"
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
	for range wavs {
		tpv.Println("Wave Filepath:", <-wavs)
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
