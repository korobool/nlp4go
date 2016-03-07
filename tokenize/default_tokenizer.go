package tokenize

type DefaultTokenizer struct {
	sentTokenizer SentenceTokenizer
	wordTokenizer Tokenizer
}

func NewDefaultTokenizer() (*DefaultTokenizer, error) {

	sentTokenizer, err := NewPunktTokenizer()
	if err != nil {
		return nil, err
	}
	wordTokenizer := NewTBWordTokenizer(true, true, nil)
	return &DefaultTokenizer{
		sentTokenizer: sentTokenizer,
		wordTokenizer: wordTokenizer,
	}, nil
}

func (t *DefaultTokenizer) Tokenize(text string) []*Token {
	tokens := []*Token{}

	sentences := t.sentTokenizer.Tokenize(text)
	for _, s := range sentences {
		_tokens := t.wordTokenizer.Tokenize(s.Text)
		for _, tok := range _tokens {
			tok.Pos += s.PosStart
			tokens = append(tokens, tok)
		}
	}
	return tokens
}
