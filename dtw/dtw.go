package dtw

import (
	"fmt"

	"github.com/r9y9/gossp/dtw"
)

// sequence1 is template, sequence2 is aligned with sequence1
func RunDtw(sequence1 [][]float64, sequence2 [][]float64, sequence1Offset int) int {
	dtwObject := dtw.DTW{ForwardStep: 10, BackwardStep: 10}

	if false {
		// Big issue here:
		// the following lines crash. We don't know why as the signal processing
		// library doesn't provide documentation on the shape of data
		dtwObject.SetTemplate(sequence1)
		path := dtwObject.DTW(sequence2)

		fmt.Println("Path ", path)

		finalIndex := path[len(path)-1]
		return sequence1Offset + finalIndex
	}

	return sequence1Offset
}
