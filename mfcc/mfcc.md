# MFCC Documentation - go-aeneas

## MFCC (Mel-Frequency Cepstral Coefficients) Overview:

- go-aeneas implementation uses a series of general-use audio-processing algorithms to isolate and extract audio features in the form of coefficients
- Sequential, pipelined audio processing comprised of six primary steps: Framing, Hamming Window, FFT, Triangular Filtering, Logarithmic Smoothing application, DCT (Discrete Cosine Transforming)
- Implemented plot-visualization functionality - available contingent upon the --plot flag. Generates a PNG with a time-series plot of the maximum coefficients at each instance in time

## Libraries Included:
#### GoNum - Plot Libraries for graphical data visualization:
- gonum.org/v1/plot
- gonum.org/v1/plot/plotter
- gonum.org/v1/plot/plotutil
- gonum.org/v1/plot/vg
#### Go-Audio for wav file processing/manipulation
- github.com/go-audio/wav
#### GoNum Fourier and Window transformation libraries for framing/transformation
- gonum.org/v1/gonum/dsp/fourier
- gonum.org/v1/gonum/dsp/window

## MFCC Pipeline
####  Load audio file
- Loads file and error checks to make sure the file opens correctlly

####  Apply Pre-emphasis
- Create an array of float64 values
- Multiply the values by the pre-emphasis coefficient (0.95)

####  Amplify the high frequencies
- Find the biggest value
- If the values are less than 0, multiply by -1 to make absolute value
- Divide all the values in the array by the maximum value

####  Frame the Signal
- Declare our constant framing parameters: sample rate (44100), frame length (0.03), frame overlap (0.5)
- Generate frame size, step, and number of frames based on sample rate, frame length and passed in signal
- Create a two-dimensional array & fill it based on the frame size and step using passed in signal

####  Create overlapping segments
- In example 1234, 12 is first frame, 23 is second frame, 34 is third frame
- Overlapping reduces artifacting
- Use a 50% frame overlap using the declared overlap coefficient

####  Apply a Window Function - Like hamming
- Copy the frames into a new array & append
- Call the hamming function on every index of the modified array

####  Apply FFT to each windowed frame
- FFT is Fast Fourier Transform: Computes the Discrete Fourier Transform (DFT), converting the signal into our frequency domain
- Converts from time domain to frequency domain
- Call the fourier.NewFFT through the entirety of the framed signal arry
- Store the fft coefficients and return

#### Create a Power Spectrum
- Clone the array of FFT Coefficients
- Compute the magnitude of each nested array index (2nd dimension) 
- Fill the power spectrum with squared magnitude of the fft

####  Apply Triangular Filtering - A type of mel filter, so still spaced on the mel scale.
- The filter banks (triangular filtering) divide the signal's frequency spectrum into multiple frequency bands so that each band can be analyzed seperately.
- Triangular Filtering is applied vertically
- Map the power spectrum onto the mel Scale using implemented algorithm

####  Apply a logarithm to the filter bank energies (Weighted Output)
-  Filter bank energies - computed by summing the magnitudes of the FFt output within each of the triangular filters
- 20 'Zero' values padded on either end of the power spectrum for the kernel
- Establish a triangular kernel array
- Multiply the padded spectrum by the kernel constant array to traingulate the spectrum

####  Take the DCT - Discrete cosine transform
- Take the DCT of the list of mel log powers
- Using dct.Transform third-party function to compute the Discrete Cosine Transforms
- The MFCCs are the amplitudes of the resulting DCT spectrum

## MFCC Visualization:
- Using GoNum/Plot and related libraries to generate a time-series graphical representation of the finalized MFCC spectrum
- Graph can be generated as a PNG (Stored in go-aeneas/mfcc) through use of the conditional '--plot' flag
- Looping through dimension 1 of the spectrum (i) and computing the maximum value of the 2nd dimension (j) for each. Maximums are plotted on the Y-Axis with the frame-index mapped on the X-Axis
