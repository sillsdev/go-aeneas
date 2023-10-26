package main

import (
	"context"
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

type contextKey int

const (
	taskKey contextKey = iota
	logBuilderKey
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

	doneCh := make(chan context.Context, len(tasks))

	for _, task := range tasks {
		taskCtx, cancel := context.WithCancel(context.Background())
		taskCtx = context.WithValue(taskCtx, taskKey, task)
		taskCtx = context.WithValue(taskCtx, logBuilderKey, &strings.Builder{})

		logHeader(taskCtx)

		go func(ctx context.Context) {
			defer wg.Done()
			defer cancel()

			runPipeline(ctx)

			doneCh <- ctx
		}(taskCtx)

		wg.Add(1)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for taskCtx := range doneCh {
		builder := taskCtx.Value(logBuilderKey).(*strings.Builder)
		fmt.Println(builder.String())
	}
}
