// +build ignore

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/korobool/nlp4go/pos"
	"github.com/korobool/nlp4go/tokenize"
	"log"
	"os"
)

var ModelPath string

func parseFlags() {
	flag.StringVar(&ModelPath, "model", "test-model.gob", "path to model file")
	flag.Parse()
}

func main() {

	parseFlags()

	scanner := bufio.NewScanner(os.Stdin)

	tokenizer := tokenize.NewTBWordTokenizer(true, true, nil)

	cfgTagger := pos.TaggerConfig{
		Tokenizer: tokenizer,
		LoadModel: true,
		ModelPath: ModelPath,
	}
	posTagger, err := pos.NewPerceptronTagger(cfgTagger)
	if err != nil {
		log.Fatalf("Failed to create POS tagger: %v", err)
	}

	for scanner.Scan() {
		line := scanner.Text()

		tokens, err := posTagger.Tag(line)
		if err != nil {
			log.Fatalf("Failed to get POS: %v", err)
		}
		for _, token := range tokens {
			fmt.Println(token)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading standard input:", err)
	}
}
