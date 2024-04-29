package audiogenerators

import (
	"fmt"
	"io"
	"os"

	"github.com/sillsdev/go-aeneas/datatypes"
)

type AudioFileCopy struct {
}

// An audio "generator" which doesn't actually generate audio but simply copies it from a different folder
// In order to work, a parameter is expected to be provided, `espeak_output_directory`, which contains all the .wav files
// with the basename being the phrase index (e.g., 1, 2a) as specified in the phrase input file
func (afc AudioFileCopy) GenerateAudioFile(parameters *datatypes.Parameters, phrase *datatypes.Phrase, outputPath string) error {
	source, err := os.Open(fmt.Sprintf("%s/%s.wav", parameters.Get("espeak_output_directory"), phrase.PhraseIndex))
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Open(outputPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func (afc AudioFileCopy) GetName() string {
	return "copy"
}

func GetAudioCopier() AudioFileCopy {
	return AudioFileCopy{}
}
