package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
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

	var wg sync.WaitGroup

	doneCh := make(chan *strings.Builder, len(tasks))

	for _, task := range tasks {
		builder := &strings.Builder{}

		logHeader(task, builder)

		go func(task *Task, builder *strings.Builder) {
			defer wg.Done()

			runPipeline(task, builder)

			doneCh <- builder
		}(task, builder)

		wg.Add(1)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for builder := range doneCh {
		fmt.Println(builder.String())
	}
}
