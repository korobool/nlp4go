package gonlp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitTokenizer(t *testing.T) {
	sentence := "Test sentence tokenization."
	tokenizer := NewSplitTokenizer(" ")
	tokens := tokenizer.Tokenize(sentence)
	assert.Equal(t, tokens[1].Word, "sentence")
}

func TestRegexpTokenizer(t *testing.T) {
	sentence := "Test sentence tokenization."
	tokenizer := NewRegexpTokenizer(`\s+`)
	tokens := tokenizer.Tokenize(sentence)
	assert.Equal(t, tokens[1].Word, "sentence")
}

func TestTreeBankTokenizer(t *testing.T) {
	sentence := `They'll save 9$ and 10% invest.`
	tokenizer := NewTreeBankTokenizer()
	tokens := tokenizer.Tokenize(sentence)
	assert.Equal(t, tokens[4].Word, "$")
}
