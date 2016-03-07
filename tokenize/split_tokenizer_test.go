package tokenize

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitSpacesASCII(t *testing.T) {
	example := "Once we accept our limits, we go beyond them."
	tokenizer := NewSplitTokenizer(" ")
	tokens := tokenizer.Tokenize(example)
	poses := []int{0, 5, 8, 15, 19, 27, 30, 33, 40}
	words := []string{"Once", "we", "accept", "our", "limits,", "we", "go", "beyond", "them."}

	assert.Equal(t, len(tokens), len(words))

	for i, token := range tokens {
		assert.Equal(t, token.Pos, poses[i])
		assert.Equal(t, token.Word, words[i])
	}
}

func TestSplitSpacesUnicode(t *testing.T) {
	example := "Видимость работы это еще не работа."
	tokenizer := NewSplitTokenizer(" ")
	tokens := tokenizer.Tokenize(example)
	poses := []int{0, 10, 17, 21, 25, 28}
	words := []string{"Видимость", "работы", "это", "еще", "не", "работа."}

	assert.Equal(t, len(tokens), len(words))

	for i, token := range tokens {
		assert.Equal(t, token.Pos, poses[i])
		assert.Equal(t, token.Word, words[i])
	}
}

func TestSplitCustomString(t *testing.T) {
	example := "Мишка%%%очень%%%любит%%%мёд."
	tokenizer := NewSplitTokenizer("%%%")
	tokens := tokenizer.Tokenize(example)
	poses := []int{0, 8, 16, 24}
	words := []string{"Мишка", "очень", "любит", "мёд."}

	assert.Equal(t, len(tokens), len(words))

	for i, token := range tokens {
		assert.Equal(t, token.Pos, poses[i])
		assert.Equal(t, token.Word, words[i])
	}
}
