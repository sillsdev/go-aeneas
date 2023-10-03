Currently existing python libraries:

ESpeak
Python speech-synthesizer program which accepts text input and generates an audio file (either mp3 or wav). Multiple different voice synthesizer options.
Djangulo go-espeak is a potential alternative: https://pkg.go.dev/github.com/djangulo/go-espeak#section-readme
Eventually we would like an eSpeak implementation as the default option for TTS, because of its local nature. Other TTS options will be added in the future for additional functionality and efficiency. These options would interact with a web-server, not storing data locally.
Another more attractive alternative is Espeak-NG: https://github.com/BenLubar/espeak
Espeak-NG uses the format synthesis method to allow for hundreds of languages to be used with small, lightweight footprint. Not as clear or as smooth as larger synthesizers. Produces output as a .wav file, features different voices and characteristics which can be changed. Written in C. Can read from a file or stdin to generate Text-to-Speech. SSML and HTML are both supported. Total footprint is only a few MB.

NumPy
Allows for complex mathematical expressions and equations. N-dimensional arrays.
Looks like it is currently being used for n-dimensional array creation and management used for generating timing files (arrays within arrays).
May be inefficient or unnecessary, may not need a third-party library to implement similar arrays. GOLang may have the internal capabilities to handle similar algorithms.
Gonum is a potential alternative: https://github.com/gonum/gonum if a third-party library is necessary.

Mfcc.py
Mel-Frequency Cepstral Coefficients being used for wave-form processing and speech pattern recognition
Within existing Aeneas infrastructure it is being used for wave-from processing and designing the timing files using pattern-recognition for phrasing.
Currently using an external C implementation which we may be able to carry over into our GOLang implementation.
No equivalent GO third-party library for the existing python library.

DTWSpeech - dtw.py
DTW refers to dynamic time warping, an algorithm used to compare the timing sequences of two similar audio waves. This algorithm uses detailed analysis and comparison of two piece-wise functions to achieve temporal alignment. 
Aeneas currently uses aeneas.cdtw Python C Extension to compute the DTW (Best path between two audio waves given the MFCC), the approximated cost matrix, and the accumulated cost matrix. All functions return tuples in an array format, should utilize slices in GO transposition.
Potential alternative is Mjanda DTW for GO: https://github.com/mjanda/go-dtw which is ported from a detailed javascript implementation of DTW.