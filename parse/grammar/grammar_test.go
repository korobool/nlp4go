package grammar

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	input = `
		S -> NP VP
		PP -> P NP
		NP -> Det N | NP PP
		VP -> V NP | VP PP
		Det -> 'a' | 'the'
		N -> 'dog' | 'cat'
		V -> 'chased' | 'sat'
		P -> 'on' | 'in'
	`

	input_pr = `
		S -> NP VP [1.0]
	     	NP -> Det N [0.5] | NP PP [0.25] | 'John' [0.1] | 'I' [0.15]
	        Det -> 'the' [0.8] | 'my' [0.2]
		N -> 'man' [0.5] | 'telescope' [0.5]
		VP -> VP PP [0.1] | V NP [0.7] | V [0.2]
		V -> 'ate' [0.35] | 'saw' [0.65]
		PP -> P NP [1.0]
		P -> 'with' [0.61] | 'under' [0.39]
	`
)

func BenchmarkGetCFG(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetContextFreeGrammarFromString(input, false)
	}
}

func BenchmarkGetCFGP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetContextFreeGrammarFromString(input_pr, true)
	}
}

func TestBasic(t *testing.T) {
	cfg, err := GetContextFreeGrammarFromString(input, false)
	assert.Equal(t, err, nil)
	assert.Equal(t, cfg.Start().Val, "S")
	plist := cfg.Productions(nil, nil, false)
	assert.Equal(t, len(plist), 14)
	assert.Equal(t, plist[0].Lhs.Val, "S")
	assert.Equal(t, len(plist[0].Rhs), 2)
	assert.Equal(t, plist[0].Rhs[0].Val, "NP")
	assert.Equal(t, plist[0].Rhs[1].Val, "VP")
	assert.Equal(t, plist[0].Rhs[1].NonTerminal, true)

	assert.Equal(t, plist[13].Lhs.Val, "P")
	assert.Equal(t, len(plist[13].Rhs), 1)
	assert.Equal(t, plist[13].Rhs[0].Val, "in")
	assert.Equal(t, plist[13].Rhs[0].NonTerminal, false)
	fmt.Println("")

	cfg, err = GetContextFreeGrammarFromString(input_pr, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, cfg.Start().Val, "S")
	plist = cfg.Productions(nil, nil, false)
	assert.Equal(t, len(plist), 17)
	assert.Equal(t, plist[16].Lhs.Val, "P")
	assert.Equal(t, len(plist[16].Rhs), 1)
	assert.Equal(t, plist[16].Rhs[0].Val, "under")
	assert.Equal(t, plist[16].Rhs[0].NonTerminal, false)
	assert.Equal(t, plist[16].Probability, 0.39)
}
