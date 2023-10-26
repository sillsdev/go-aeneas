package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	logLevel    = 0
	batch       = ""
	logBuilders *LogBuilders
)

type contextKey int

const (
	taskIdKey contextKey = iota
	taskKey
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

	doneCh := make(chan int, len(tasks))
	logBuilders = NewLogBuilders(len(tasks))

	for taskId, task := range tasks {
		ctx, cancel := context.WithCancel(context.Background())

		taskCtx := context.WithValue(ctx, taskIdKey, taskId)
		taskCtx = context.WithValue(taskCtx, taskKey, task)

		logHeader(taskId, task)

		go func(ctx context.Context) {
			defer wg.Done()
			defer cancel()

			runPipeline(ctx)

			taskId := ctx.Value(taskIdKey).(int)
			doneCh <- taskId
		}(taskCtx)

		wg.Add(1)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for taskId := range doneCh {
		builder := logBuilders.Get(taskId)
		fmt.Println(builder.String())
	}
}
