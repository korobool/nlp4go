package tokenize

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSentenceSplit(t *testing.T) {
	example := "Julia with half an ear listened to the list Margery read out and, though she knew the room so well, idly looked about her. It was a very proper room for the manager of a first-class theatre. The walls had been panelled (at cost price) by a good decorator and on them hung engravings of theatrical pictures by Zoffany and de Wilde. The armchairs were large and comfortable."

	tokenizer, err := NewPunktTokenizer()
	if err != nil {
		t.Fail()
	}
	tokens := tokenizer.Tokenize(example)
	posStarts := []int{0, 123, 191, 331}
	posEnds := []int{121, 189, 329, 371}

	assert.Equal(t, len(tokens), len(posStarts))

	for i, token := range tokens {
		t.Logf("Text:>>%s<<", token.Text)
		assert.Equal(t, posStarts[i], token.PosStart)
		assert.Equal(t, posEnds[i], token.PosEnd)
	}
}
