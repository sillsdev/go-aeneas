# Go Aeneas Handoff

## MFCC Known Issues

mfcc.go - func mfccWeighSignal(powerSpectrum [][]float64) [][]float64
- The use of copy may be causing a proliferation of null values at the start and end of the power spectrum.

## DTW implementation crashes with multi-dimensional inputs

Shape of data requirements is not provided via documentation, and directly providing data from MFCC causes crashes 

## Additionally, ESpeakNG does not compile or run outside of a Linux Environment

More information provided in the espeak-wasm branch