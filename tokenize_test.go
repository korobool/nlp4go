package gonlp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitTokenizer(t *testing.T) {
	sentence := "Test sentence tokenization."
	tokenizer := NewSplitTokenizer(" ")
	tokens := tokenizer.Tokenize(sentence)
	assert.Equal(t, "sentence", tokens[1].Word)
}

func TestRegexpTokenizer(t *testing.T) {
	sentence := "Test sentence tokenization."
	tokenizer := NewRegexpTokenizer(`\s+`)
	tokens := tokenizer.Tokenize(sentence)
	assert.Equal(t, "sentence", tokens[1].Word)
}

func TestTreeBankTokenizer(t *testing.T) {
	sentence := `"They'll" save 9$ and 10% we'll invest.`
	tokenizer := NewTreeBankTokenizer()
	tokens := tokenizer.Tokenize(sentence)
	fmt.Println(tokens)
	assert.Equal(t, "``", tokens[0].Word)
	assert.Equal(t, "''", tokens[3].Word)
	assert.Equal(t, "$", tokens[6].Word)
}
