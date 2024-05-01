package datatypes

type AudioGenerator interface {
	GenerateAudioFile(parameters *Parameters, phrase *Phrase, outputPath string) error
	GetName() string
}
