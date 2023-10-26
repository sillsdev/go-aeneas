package main

import (
	"fmt"
	"strings"
)

func logHeader(task *Task, builder *strings.Builder) {
	if len(task.Description) > 0 {
		builder.WriteString(fmt.Sprintln(""))
		builder.WriteString(fmt.Sprintln("*** ", task.Description, " ***"))
		builder.WriteString(fmt.Sprintln(""))
	}

	parameters := parseParameters(task.Parameters)

	builder.WriteString(fmt.Sprintln("Audio   : ", task.AudioFilename))
	builder.WriteString(fmt.Sprintln("Phrase  : ", task.PhraseFilename))
	builder.WriteString(fmt.Sprintln("Output  : ", task.OutputFilename))
	builder.WriteString(fmt.Sprintln("Parameters : ", parameters))
}
