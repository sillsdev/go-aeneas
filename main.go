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

	for _, task := range tasks {
		if len(task.Description) > 0 {
			fmt.Println("")
			fmt.Println("*** ", task.Description, " ***")
			fmt.Println("")
		}

		parameters := parseParameters(task.Parameters)

		fmt.Println("Audio   : ", task.AudioFilename)
		fmt.Println("Phrase  : ", task.PhraseFilename)
		fmt.Println("Output  : ", task.OutputFilename)
		fmt.Println("Parameters : ", parameters)
	}
}
