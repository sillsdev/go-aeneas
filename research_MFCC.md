MFCC Pipeline and Outline + Libraries and Functions

__Overview__
1. Load audio file
    Retrieve the signal. Load into an audio buffer?
    - Decode the signal

2. Apply Pre-emphasis
    Amplify the high frequencies - protects against noise
	y[n] = x[n] - α*x[n-1]
		y[n] is the output or emphasized signal
		x[n] is the input signal
		a is the pre-emphasis coefficient, usually between 0.9 to 1
			// depending on the use case.
3. Frame the Signal
        Create overlapping segments (overlapping reduces artifacting?)
4. Apply a Window Function - Like hamming
        Each frame gets "windowed" by multiplying it with a window function
	    Overlapping windows is almost required. Beneficial when used DTW
	    and/because it preserves temporal information
5. Apply FFT to each windowed frame
        Converts from time domain to frequency domain
6. Apply Triangular Filtering - A type of mel filter, so still spaced on the mel scale.
        The filter banks (triangular filtering) divide the signal's frequency spectrum into multiple frequency bands.
        So that each band can be analyzed separately.
7. Apply a logarithm to the filter bank energies
        Filter bank energies - computed by summing the magnitudes of the FFt output within each of the triangular filters
8. Take the DCT - Discrete cosine transform
        Take the DCT of the list of mel log powers
        The MFCCs are the amplitudes of the resulting spectrum




Go Libraries/Functions to accomplish the above
__Go audio__
https://pkg.go.dev/github.com/go-audio/audio#Buffer
    1. Buffer structure. Load the entire initial signal to a buffer.
	May want to use a buffer between tasks, However since the signal
	is quite small, it may be unnecessary. Though there may be other
	reasons to use the buffer throughout.


__Go audio/wav__
https://pkg.go.dev/github.com/go-audio/wav
    1. Decoder to load the wav file and gather the signal, sample rate
            Loaded into a Format struct
            PCM format - holds the relevant parameters (Pulse Code Modulation)
	
__Gonum__
https://pkg.go.dev/gonum.org/v1f/gonum/dsp/fourier
        fft function

__io__
https://pkg.go.dev/io#ReadSeeker

__riff__ -- May not be needed manually
https://pkg.go.dev/github.com/go-audio/riff#section-readme

__mel__
https://pkg.go.dev/github.com/emer/auditory/mel


__No library needed for these?__
    2. Pre-emphasis
        y[n] = x[n] - α*x[n-1] : this is the function used to apply pre-emphasis
		https://tinyurl.com/Pre-emphasis-Stack-Overflow
	    3. Framing
		    Manual build using gonum
	    4. Windowing
		    Manual build using gonum


__Step Through__

1. Loading and Decoding the Audio file
    - It is loaded using a decoder with the PCM Format (Either FullPCMBuffer or PCMBuffer)
        - FullPCMBuffer loads the entirety into memory and is indicated to be inefficient.
            - However, due to our small audio file size this may not be an issue.
        - PCMBUffer may be the preffered choice overall

Func getSamples(string pathToAudioFile) and return a slice of samples? (array=fixed, slice=dynamic){  
    
    myAudioFile = open(pathToAudioFile)
        // Do an error check to see if file exists/opened
        
    // Decode the audio file
        // Relevant documentation - https://github.com/go-audio/wav/blob/v1.1.0/decoder.go#L238
            // NewDecoder - Decoder - PCMBuffer
        // NewDecoder uses io.ReadSeeker which should work with wav files.


    myDecoder = NewDecoder(myAudioFile)

    // func (d *Decoder) PCMBuffer(buf *audio.IntBuffer) (n int, err error) {
    // So a buffer needs to be created and passed.

    audioBuffer = audio.IntBuffer // Make this a pointer // Int version of the PCM data

    myDecoder.PCMBuffer(audioBuffer) // Pass pointer here


    // extract the samples from the buffer - converting to float to prepare for pre-emphasis
    // https://github.com/go-audio/audio/blob/v1.0.0/int_buffer.go#L15
        // Use sourcebitdepth in a calculation to normalize and convert to float.

    float64[of size audioBuffer] samples = audioBuffer converted to float and normalized


    return samples

}

// Referrence old-aeneas to see what alpha value they used.

func preEmphasis(float64[]samples, the alpha value between 0.9 and 1.0 we want to use ){
    // Take in the slice of samples
    // Apply the pre-emphasis calculation

    float64[of size samples] filteredSignal = samples[i] - alpha*samples[i-1]

    return filteredSignal
}

// Other functions still need to be outlined.
