package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	logLevel = 0
	batch    = ""
)

func processTask(results chan string, task *Task) {
	logs := strings.Builder{}

	if len(task.Description) > 0 {
		logs.WriteString(fmt.Sprintln(""))
		logs.WriteString(fmt.Sprintln("*** ", task.Description, " ***"))
		logs.WriteString(fmt.Sprintln(""))
	}

	parameters := parseParameters(task.Parameters)

	logs.WriteString(fmt.Sprintln("Audio   : ", task.AudioFilename))
	logs.WriteString(fmt.Sprintln("Phrase  : ", task.PhraseFilename))
	logs.WriteString(fmt.Sprintln("Output  : ", task.OutputFilename))
	logs.WriteString(fmt.Sprintln("Parameters : ", parameters))

	results <- logs.String()
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
