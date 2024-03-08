package datatypes

import (
	"fmt"
	"path/filepath"
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
	Generator  *AudioGenerator
	TempDir    string
}

func NewTaskProcessVariables(task *Task, generator *AudioGenerator, tempDir string) *TaskProcessVariables {
	return &TaskProcessVariables{
		Task:       task,
		Logs:       strings.Builder{},
		Parameters: ParseParameters(task.Parameters),
		Generator:  generator,
		TempDir:    tempDir,
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

func (tpv *TaskProcessVariables) GetWavFilepath() string {
	return filepath.Join(tpv.TempDir, filepath.Base(tpv.Task.AudioFilename)+".wav")
}

func (tpv *TaskProcessVariables) GetPhraseFilePath(phraseIndex string) string {
	return filepath.Join(tpv.TempDir, fmt.Sprintf("%s.%s.wav", filepath.Base(tpv.Task.AudioFilename), phraseIndex))
}
