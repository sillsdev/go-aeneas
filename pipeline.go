package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func runPipeline(ctx context.Context) {
	task := ctx.Value(taskKey).(*Task)
	taskId := ctx.Value(taskIdKey).(int)

	builder := logBuilders.Get(taskId)

	builder.WriteString(fmt.Sprintln("Running pipeline..."))

	// Do work here
	builder.WriteString(fmt.Sprintf("Doing Task: %s\n", task.Description))
	sleepDuration := time.Duration(rand.Intn(10)+1) * time.Second
	time.Sleep(sleepDuration)

	builder.WriteString(fmt.Sprintln("Pipeline finished."))
}
