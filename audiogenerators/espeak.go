package audiogenerators

import (
	"fmt"
	"os"

	"github.com/sillsdev/go-aeneas/datatypes"

	"gopkg.in/BenLubar/espeak.v2"
)

// C Parent: https://github.com/espeak-ng/espeak-ng/blob/master/docs/guide.md
// https://github.com/readbeyond/aeneas/blob/master/aeneas/synthesizer.py

// function inputs the task struct and outputs a .wav file
type EspeakGenerator struct {
	ctx espeak.Context
}

func (gen EspeakGenerator) GenerateAudioFile(parameters *datatypes.Parameters, phrase string, outputPath string) error {
	language := parameters.Get("language")

	//similar to printf in C, prints to the string
	//the %s gets replaced with the passed in parameters
	phrase_ssml := fmt.Sprintf(`
		<?xml version="1.0"?>
		<speak version="1.1"
			xmlns="http://www.w3.org/2001/10/synthesis"
			xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
			xsi:schemaLocation="http://www.w3.org/2001/10/synthesis
				http://www.w3.org/TR/speech-synthesis11/synthesis.xsd"
			xml:lang="en-US">
			<voice gender="male" languages="%s">
				"%s"
			</voice>
		</speak>
	`, language, phrase)

	ctx := gen.ctx
	ctx.SynthesizeText(phrase_ssml)

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	ctx.WriteTo(f)

	return nil
}

func (gen EspeakGenerator) GetName() string {
	return "espeak-ng"
}

func GetEspeakGenerator() *EspeakGenerator {
	return &EspeakGenerator{ctx: espeak.Context{}}
}
