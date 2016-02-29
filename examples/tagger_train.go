package main

import (
	"flag"
	"fmt"
	"github.com/korobool/nlp4go/pos"
	"github.com/korobool/nlp4go/tokenize"
	"github.com/korobool/nlp4go/utils"
	"log"
)

var (
	OntonotesPath string
	ModelPath     string
)

func parseFlags() {
	flag.StringVar(&OntonotesPath, "corpus", "ontonotes", "path to OntoNotes dir")
	flag.StringVar(&ModelPath, "model", "test-model.gob", "path to model file")
	flag.Parse()
}

func main() {

	parseFlags()

	parser, err := utils.NewOntonotesParser(OntonotesPath)
	if err != nil {
		log.Fatalf("Failed to create parser: %v", err)
	}
	trainSentences, err := parser.ParseToWordsTags()
	if err != nil {
		log.Fatalf("Failed to parse OntoNotes: %v", err)
	}

	tokenizer := tokenize.NewTBWordTokenizer(true, true, nil)

	cfgTagger := pos.TaggerConfig{
		Tokenizer: tokenizer,
		LoadModel: false,
		ModelPath: ModelPath,
	}
	posTagger, err := pos.NewPerceptronTagger(cfgTagger)
	if err != nil {
		log.Fatalf("Failed to create POS tagger: %v", err)
	}

	var progress int
	progressCallback := func(current, total int) {
		p := (100 * current) / total
		if p > progress {
			progress = p
			fmt.Printf("\rTrain:%4d%%", progress)
		}
	}

	posTagger.Train(trainSentences, 5, progressCallback)
	posTagger.SaveModel()
}
