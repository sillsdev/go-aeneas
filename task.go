package main

type Task struct {
	Description    string `json:"description"`
	AudioFilename  string `json:"audioFilename"`
	PhraseFilename string `json:"phraseFilename"`
	Parameters     string `json:"parameters"`
	OutputFilename string `json:"outputFilename"`
}
