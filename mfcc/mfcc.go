package mfcc

import (
	"fmt"
	"math"
	"os"

	"github.com/go-audio/wav"
	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/dsp/window"
)

func GenerateMfcc(inFileName string) ([][]float64, error) {
	fmt.Println("Beginning mfcc generation for inputted .wav file: ", inFileName)

	signal, err := mfccLoadSignal(inFileName)
	if err != nil {
		return nil, err
	}
	normalizedSignal := mfccNormalize(signal)
	framedSignal := mfccFrameSignal(normalizedSignal)
	windowedSignal := mfccWindowSignal(framedSignal)
	fft := mfccFFT(framedSignal, windowedSignal)
	powerSpectrum := mfccPowerSpectrum(fft)
	triangularFiler := mfccTriangularFilter(powerSpectrum)
	weightedSignal := mfccWeighSignal(triangularFiler)
	mfcc := mfccDCT(weightedSignal)

	return mfcc, nil
}

func mfccLoadSignal(inFileName string) ([]float64, error) {

	const preEmphasis = 0.95 // PreEmphasis Coefficient -> Modify coefficient as needed

	audiofile, err := os.Open(inFileName) // OS opens inFileName; takes a string filepath
	if err != nil {
		return nil, err
	}

	decoder := wav.NewDecoder(audiofile) // Load a decoder for the loaded wav file
	if decoder == nil {
		return nil, fmt.Errorf("could not decode file")
	}

	audiobuffer, err := decoder.FullPCMBuffer() // FullPCMBuffer takes a pointer to a decoder and returns a buffer
	if err != nil {
		return nil, err
	}

	defer audiofile.Close()

	// End Signal Loading & Apply Preemphasis

	signal := audiobuffer.Data

	//Load audio buffer data into new signal array
	signal64 := make([]float64, len(signal))
	for i := 0; i < len(signal); i++ {
		signal64[i] = float64(signal[i])
	}
	for i := 1; i < len(signal); i++ {
		signal64[i] = signal64[i] - preEmphasis*signal64[i-1]
	}

	return signal64, nil
}

func mfccNormalize(signal64 []float64) []float64 {
	currentMax := signal64[0]
	if currentMax < 0 {
		currentMax = currentMax * -1
	}

	for i := 1; i < len(signal64); i++ {
		next := signal64[i]
		if next < 0 {
			next = next * -1
		}
		if next > currentMax {
			currentMax = next
		}
	}
	absMax := currentMax

	for i := 0; i < len(signal64); i++ {
		signal64[i] = signal64[i] / absMax
	}

	return signal64
}

// 44100(samplingRate) * 0.03(milliseconds) = 1323 (frame size)
// Sampling Rate = 44.1kHz
// Frame Length = 20-30ms is a good choice

func mfccFrameSignal(signal64 []float64) [][]float64 {
	var sampleRate = 44100.0
	var frameLength = 0.03
	var frameOverlap = 0.5
	var frameSize = int(sampleRate * frameLength)                                        // int: can't have a fraction of a sample
	var frameStep = int(sampleRate * frameLength * frameOverlap)                         // int: indexes are whole numbers
	var numFrames = (float64(len(signal64)) / (sampleRate * frameLength)) / frameOverlap // +1 if partial
	if numFrames != float64(int(numFrames)) {                                            // If there is a partial frame truncate and add 1, then handle the the partial/tail frame.
		numFrames = float64(int(numFrames))
	}
	fmt.Println("signal64: ", len(signal64))
	fmt.Println("frameSize: ", frameSize)
	fmt.Println("frameStep: ", frameStep)
	fmt.Println("numFrames: ", numFrames)

	frames := make([][]float64, int(numFrames))
	for i := 0; i < int(numFrames); i++ {
		start := i * frameStep
		end := start + frameSize
		if end > len(signal64) {
			end = len(signal64)
		}
		frames[i] = signal64[start:end]
	}

	return frames
}

// End Framing the Signal
// Start Windowing the Signal (Hamming Window)
// https://pkg.go.dev/gonum.org/v1/gonum/dsp/window#example-Hamming

// Have to copy the frames because I need both the original and windowed frames for the FFT.

func mfccWindowSignal(frames [][]float64) [][]float64 {
	goingHam := make([][]float64, len(frames))
	for i := 0; i < len(frames); i++ {
		goingHam[i] = append(goingHam[i], frames[i]...)
	}

	for i := 0; i < len(goingHam); i++ {
		window.Hamming(goingHam[i]) // Changes data in place according to documentation
	}

	return goingHam
}

// Now the overlap-add step. https://en.wikipedia.org/wiki/Overlap%E2%80%93add_method
// recombining the frames accounting for the overlap.

// - FFT then overlap-add
// - Overlap-add method
// There doesn't seem to be a library for this, so, like framing, do it manually.
// So even though the slice of slices are independent from each other and so there isn't an issue with the windowing function changing the values of the other frames, I still need to account for the overlap?

func mfccFFT(frames [][]float64, goingHam [][]float64) [][]complex128 {
	fftCoefficients := make([][]complex128, len(frames))

	for i := 0; i < len(frames); i++ {
		fft := fourier.NewFFT(len(frames[i]))
		fftCoefficients[i] = fft.Coefficients(nil, goingHam[i])
	}

	return fftCoefficients
}

// https://en.wikipedia.org/wiki/Triangular_function
// Convolution: Operation that combines two signals to produce a third signal. - this should include/be the overlap-add step.
// https://pkg.go.dev/github.com/brettbuddin/fourier@v0.1.1#section-readme
// END FFT

// Start Power Spectrum - Triangular filtering Pre-Step
// I'm squaring the magnitude (so abs value) of the complex numbers, But im extracting the real numbers and discarding the imaginary numbers.

func mfccPowerSpectrum(fftCoefficients [][]complex128) [][]float64 {

	powerSpectrum := make([][]float64, len(fftCoefficients))
	for i := 0; i < len(fftCoefficients); i++ {
		powerSpectrum[i] = make([]float64, len(fftCoefficients[i]))
		for j := 0; j < len(fftCoefficients[i]); j++ {
			magnitude := math.Sqrt(real(fftCoefficients[i][j])*real(fftCoefficients[i][j]) + imag(fftCoefficients[i][j])*imag(fftCoefficients[i][j]))
			powerSpectrum[i][j] = magnitude * magnitude
		}
	}

	return powerSpectrum
}

// End Power Spectrum

// Start Triangular Filtering - Mel Filter Banks
// https://www.statistics.com/glossary/triangular-filter/#:~:text=As%20compared%20to%20the%20rectangular,short%20course%20Time%20Series%20Forecasting%20
// https://en.wikipedia.org/wiki/Mel_scale
// https://en.wikipedia.org/wiki/Periodogram

// Take the magnitude squared of the complex Fourier coefficients - This is the power spectrum
// Map the power spectrum onto the Mel scale using a filterbank.

// https://en.wikipedia.org/wiki/Mel_scale
//melScale := 2595 * math.Log10(1+(sampleRate/2)/700) // Divide by 2 for Nyquist frequency
//filterBankSize := 26                                // Number of filters

// NOTE! This is NOT like windowing/framing. Framing we can think of as something we applied horizontally all along the data. The filter(s) we are applying we apply vertically.

// Map the power spectrum onto the Mel scale
// melScale := 2595 * math.Log10(1+(each_index_of_powerSpectrum/2)/700) // Divide by 2 for Nyquist frequency

func mfccTriangularFilter(powerSpectrum [][]float64) [][]float64 {
	for i := 0; i < len(powerSpectrum); i++ {
		for j := 0; j < len(powerSpectrum[i]); j++ {
			powerSpectrum[i][j] = 2595 * math.Log10(1+(float64(powerSpectrum[i][j])/2)/700) // Divide by 2 for Nyquist frequency
		}
	}
	return powerSpectrum
}

// ---Filter Bank Process---
// The filter banks are vertical
// Apply the filter banks to the power spectrum
// For range of PowerSpectrum multiply x adjacent indices by the filter bank values
// add them together to get the weighted sum? take the log of the result

// Does every filterbank get applied across the entire spectrum? or is it dynamic, one filter per "section size"?
// Meaning the entire set of filterbanks applied to every part of the spectrum.
// End Triangular Filtering
// [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] - indexes
// [1, 1, 1, 1, 1, 1, 1, 1, 1, 1] - values

// Weighted output through summing
// [0+12, 0+1+23, 01+2+34, 12+3+45, 23+4+56, 34+5+67, 45+6+78, 56+7+89, 67+8+9, 78+9]

// [0.000 0.053 0.105 0.158 0.211 0.263 0.316 0.368 0.421 0.474 0.526 0.579 0.632 0.684 0.737 0.789 0.842 0.895 0.947 1.000, 0.947,]

// Add 20 zeros to the beginning and end of the power spectrum as padding for the kernel
func mfccWeighSignal(powerSpectrum [][]float64) [][]float64 {
	paddedSpectrum := make([][]float64, len(powerSpectrum))
	for i := 0; i < len(powerSpectrum); i++ {
		temp := make([]float64, len(powerSpectrum[i])+40)
		copy(temp[20:], powerSpectrum[i])
		paddedSpectrum[i] = temp
	}

	// My understanding is that for the kernel, if it is a kernel of width 40, the actual length is 41 so that the center is 1.0
	triangularkernel := []float64{0.05, 0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4, 0.45, 0.5, 0.55, 0.6, 0.65, 0.7, 0.75, 0.8, 0.85, 0.9, 0.95, 1.0, 0.95, 0.9, 0.85, 0.8, 0.75, 0.7, 0.65, 0.6, 0.55, 0.50, 0.45, 0.40, 0.35, 0.30, 0.25, 0.20, 0.15, 0.10, 0.05}

	weightedOutput := make([][]float64, len(powerSpectrum))
	for i := 0; i < len(powerSpectrum); i++ { // Index of slices
		weightedOutput[i] = make([]float64, len(powerSpectrum[i]))
		for j := 0; j < len(powerSpectrum[i]); j++ { // Index of values
			for k := 0; k < len(triangularkernel); k++ { // Index of kernel
				weightedOutput[i][j] += paddedSpectrum[i][j+k] * triangularkernel[k]
			}
		}
	}

	for i := 0; i < len(weightedOutput); i++ {
		for j := 0; j < len(weightedOutput[i]); j++ {
			weightedOutput[i][j] = math.Log(weightedOutput[i][j])
		}
	}
	return weightedOutput
}

func mfccDCT(weightedOutput [][]float64) [][]float64 {
	dctCoefficients := make([][]float64, len(weightedOutput))

	for i := 0; i < len(weightedOutput); i++ {
		dctCoefficients[i] = make([]float64, len(weightedOutput[i]))
		dct := fourier.NewDCT(len(weightedOutput[i]))
		dct.Transform(dctCoefficients[i], weightedOutput[i])
	}

	return dctCoefficients
}

// Important notes to revisit
// Dynamic vs Static Framing
// Discard half of the FFT? Only compute the first half of the FFT? Only the first half contains unique information?
// filterbank/more filters? Frequency bins is huh?
