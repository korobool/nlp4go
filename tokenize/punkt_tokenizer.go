package tokenize

import (
	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/english"
)

// Stub for PunktTokenizer. Temporarly implemented
// as wrapper around github.com/neurosnap/sentences libraray
type PunktTokenizer struct {
	sentTokenizer *sentences.DefaultSentenceTokenizer
}

func NewPunktTokenizer() (*PunktTokenizer, error) {

	sentTokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		return nil, err
	}
	return &PunktTokenizer{
		sentTokenizer: sentTokenizer,
	}, nil
}

func (t *PunktTokenizer) Tokenize(text string) []*Sentence {
	sentences := []*Sentence{}

	tmpSents := t.sentTokenizer.Tokenize(text)
	for _, s := range tmpSents {
		sentences = append(sentences, &Sentence{
			PosStart: s.Start,
			PosEnd:   s.End - 1,
			Text:     s.Text,
		})
	}
	return sentences
}
