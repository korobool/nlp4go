package main

import (
	"bufio"
	"fmt"
	"github.com/korobool/nlp4go/pos"
	"github.com/korobool/nlp4go/tokenize"
	"log"
	"os"
)

const (
	ModelPath = "test-model.gob"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	tokenizer := tokenize.NewTBWordTokenizer(true, true, nil)

	cfgTagger := pos.TaggerConfig{
		Tokenizer: tokenizer,
		LoadModel: true,
		ModelPath: ModelPath,
	}
	posTagger, err := pos.NewPerceptronTagger(cfgTagger)
	if err != nil {
		log.Fatalf("failed to create POS tagger: %v", err)
	}

	for scanner.Scan() {
		line := scanner.Text()

		tokens, err := posTagger.Tag(line)
		if err != nil {
			log.Fatalf("failed to get POS: %v", err)
		}
		for _, token := range tokens {
			fmt.Println(token)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
