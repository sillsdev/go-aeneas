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

## Development setup

Install Docker, VS Code, and the dev containers extension for VS Code.
Clone this repository with git and open it in VS Code
There should be a prompt in the bottom right to 'Reopen in dev container'; click the button to confirm this action
If this does not show up, press Ctrl-Shift-P to open the command pallette and search for 'Dev Containers: Reopen in Container'