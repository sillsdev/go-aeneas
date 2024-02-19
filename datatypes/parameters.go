package datatypes

import (
	"strings"
)

// https://stackoverflow.com/a/40380147
// Created a struct to embed map and include into context
type Parameters struct {
	m map[string]string
}

func (p Parameters) Get(key string) string {
	return p.m[key]
}

func ParseParameters(parameterString string) *Parameters {
	subParameters := strings.Split(parameterString, "|")
	mapParameters := make(map[string]string)

	for _, paramValue := range subParameters {
		kvParams := strings.Split(paramValue, "=")
		mapParameters[kvParams[0]] = kvParams[1]
	}

	parameters := &Parameters{mapParameters}

	return parameters
}
