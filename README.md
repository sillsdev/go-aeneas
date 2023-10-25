# go-aeneas
Project Pipeline:

- Begin with parsing parameters
- Put the parameters in a key value map
- Prepare log collection for go routines (create string buffer)
- Start go routines
- Start processTask function
    - Generate audio from text file (eSpeak)
    - Generate MFC coefficients from input and generated audio files
    - Use MFCC/DTW to compare the coefficients of the files
    - Write new timing to file
- End processTask
- Collect logs (buffer) and print to console
