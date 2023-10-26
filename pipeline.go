package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func runPipeline(task *Task, builder *strings.Builder) {
	builder.WriteString(fmt.Sprintln("Running pipeline..."))

	// Do work here
	builder.WriteString(fmt.Sprintf("Doing Task: %s\n", task.Description))
	sleepDuration := time.Duration(rand.Intn(10)+1) * time.Second
	time.Sleep(sleepDuration)

	builder.WriteString(fmt.Sprintln("Pipeline finished."))
}
