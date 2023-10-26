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
	log := make(chan string)
	done := make(chan struct{})

	go func() {
		sb := strings.Builder{}
		defer func() {
			results <- sb.String()
		}()

		for {
			select {
			case msg := <-log:
				sb.WriteString(msg)

			case <-done:
				return
			}
		}
	}()

	// part taken from main
	if len(task.Description) > 0 {
		log <- fmt.Sprintln("")
		log <- fmt.Sprintln("*** ", task.Description, " ***")
		log <- fmt.Sprintln("")
	}

	parameters := parseParameters(task.Parameters)

	log <- fmt.Sprintln("Audio   : ", task.AudioFilename)
	log <- fmt.Sprintln("Phrase  : ", task.PhraseFilename)
	log <- fmt.Sprintln("Output  : ", task.OutputFilename)
	log <- fmt.Sprintln("Parameters : ", parameters)

	done <- struct{}{}
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

	// let processing take place in the background...

	for range tasks {
		fmt.Println(<-results)
	}
}
