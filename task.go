package main

import (
	"fmt"
	"strings"
)

type Task struct {
	Description    string `json:"description"`
	AudioFilename  string `json:"audioFilename"`
	PhraseFilename string `json:"phraseFilename"`
	Parameters     string `json:"parameters"`
	OutputFilename string `json:"outputFilename"`
}

type TaskProcessVariables struct {
	Task       *Task
	Logs       strings.Builder
	Parameters *Parameters
}

func NewTaskProcessVariables(task *Task) *TaskProcessVariables {
	return &TaskProcessVariables{
		Task:       task,
		Logs:       strings.Builder{},
		Parameters: parseParameters(task.Parameters),
	}
}

func (tpv *TaskProcessVariables) Println(args ...interface{}) {
	tpv.Logs.WriteString(fmt.Sprintln(args...))
}

func (tpv *TaskProcessVariables) GetParameter(param string) string {
	return tpv.Parameters.Get(param)
}

func (tpv *TaskProcessVariables) GetFinalLogs() string {
	return tpv.Logs.String()
}
