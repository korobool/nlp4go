package tokenize

type Tokenizer interface {
	Tokenize(string) []*Token
}
