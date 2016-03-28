package syntax

import (
	G "../grammar"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	source = `S -> NP VP
	  VP -> V NP | V NP PP
	  PP -> P NP
	  V -> "saw" | "ate" | "walked"
	  NP -> "John" | "Mary" | "Bob" | Det N | Det N PP
	  Det -> "a" | "an" | "the" | "my"
	  N -> "man" | "dog" | "cat" | "telescope" | "park"
	  P -> "in" | "on" | "by" | "with"`
)

func BenchmarkShiftReduce(b *testing.B) {
	cfg, _ := G.GetContextFreeGrammarFromString(source, false)
	srp := ShiftReduceParser{grammar: cfg}
	gs := []*G.GrammarSymbol{G.Term("Mary"), G.Term("saw"), G.Term("a"), G.Term("dog")}
	for i := 0; i < b.N; i++ {
		_, _ = srp.Parse(gs, false)
	}
}

func TestBasic(t *testing.T) {
	cfg, err := G.GetContextFreeGrammarFromString(source, false)

	fmt.Println("source: ", source)
	fmt.Println("err: ", err)
	fmt.Println("Start: ", cfg.Start())
	fmt.Println("Productions: ")
	for _, prod := range cfg.Productions(nil, nil, false) {
		fmt.Println(prod)
	}

	srp := ShiftReduceParser{grammar: cfg}
	gs := []*G.GrammarSymbol{G.Term("Mary"), G.Term("saw"), G.Term("a"), G.Term("dog")}
	fmt.Println("=======================================")
	ret, err := srp.Parse(gs, true)
	fmt.Println("=======================================")
	assert.Equal(t, err, nil)
	assert.Equal(t, ret.Label(), "S")
	assert.Equal(t, ret.String(), "(S (NP Mary) (VP (V saw) (NP (Det a) (N dog))))")
}
