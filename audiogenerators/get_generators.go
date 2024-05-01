//go:build !linux

package audiogenerators

import "github.com/sillsdev/go-aeneas/datatypes"

func GetAudioGenerators() []datatypes.AudioGenerator {
	return []datatypes.AudioGenerator{GetAudioCopier()}
}
