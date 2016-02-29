package pos

import "github.com/korobool/nlp4go/tokenize"

type WordsTags struct {
	Words []string
	Tags  []string
}

type PosTagger interface {
	Tag(string) []tokenize.Token
}

type BaseTagger struct {
	Tokenizer tokenize.Tokenizer
}

type TaggerConfig struct {
	Tokenizer tokenize.Tokenizer
	LoadModel bool
	ModelPath string
}
