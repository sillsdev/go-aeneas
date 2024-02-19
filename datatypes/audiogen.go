package datatypes

type AudioGenerator interface {
	GenerateAudioFile(parameters *Parameters, phrase string, outputPath string) error
	GetName() string
}
