package tokenize

type Tokenizer interface {
	Tokenize(string) []*Token
}

type SentenceTokenizer interface {
	Tokenize(string) []*Sentence
}
