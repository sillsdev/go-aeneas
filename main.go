package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var (
	logLevel = 0
	batch    = ""
)

func processTask(results chan string, task *Task) {
	tpv := NewTaskProcessVariables(task)
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
}

func main() {
	processArguments()

	tasks := []*Task{}
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
		task := &Task{"", os.Args[1], os.Args[2], os.Args[3], os.Args[4]}
		tasks = append(tasks, task)
	}

	results := make(chan string)

	for _, task := range tasks {
		go processTask(results, task)
	}

	for range tasks {
		fmt.Println(<-results)
	}
}
