package datatypes

type AudioGenerator interface {
	GenerateAudioFile(parameters *Parameters, phrase string, outputPath string) error
	Close()
}

type AudioGeneratorFactory interface {
	GetAudioGenerator() (*AudioGenerator, error)
	GetName() string
}
