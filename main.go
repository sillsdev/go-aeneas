package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	logLevel = 0
	batch    = ""
)

func processTask(ctx context.Context, results chan string, task *Task) {
	sb := strings.Builder{}

	if len(task.Description) > 0 {
		sb.WriteString(fmt.Sprintln(""))
		sb.WriteString(fmt.Sprintln("*** ", task.Description, " ***"))
		sb.WriteString(fmt.Sprintln(""))
	}

	parameters := parseParameters(task.Parameters)

	sb.WriteString(fmt.Sprintln("Audio   : ", task.AudioFilename))
	sb.WriteString(fmt.Sprintln("Phrase  : ", task.PhraseFilename))
	sb.WriteString(fmt.Sprintln("Output  : ", task.OutputFilename))
	sb.WriteString(fmt.Sprintln("Parameters : ", parameters))

	results <- sb.String()
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

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	results := make(chan string, len(tasks))

	for _, task := range tasks {
		go processTask(ctx, results, task)
	}

	for range tasks {
		fmt.Println(<-results)
	}
}
