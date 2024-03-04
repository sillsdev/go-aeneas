package datatypes

import "strings"

type Phrase struct {
	PhraseIndex string
	PhraseText  string
}

func ParsePhrase(phraseLine string) (*Phrase, error) {
	phraseParts := strings.Split(phraseLine, "|")

	if len(phraseParts) < 2 {
		return nil, newError("Phrase line does not have enough parts")
	}

	return &Phrase{
		PhraseIndex: phraseParts[0],
		PhraseText:  phraseParts[1],
	}, nil
}

type PhraseParseError struct {
	msg string
}

func (err *PhraseParseError) Error() string {
	return err.msg
}

func newError(msg string) *PhraseParseError {
	return &PhraseParseError{msg}
}
