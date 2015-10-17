package main

import (
	"bufio"
	"fmt"
	"github.com/korobool/nlp4go"
	"os"
	"strings"
)

var TEXT = `Pierre|NNP Vinken|NNP ,|, 61|CD years|NNS old|JJ ,|, will|MD\n
join|VB the|DT board|NN as|IN a|DT nonexecutive|JJ director|NN\n
Nov.|NNP 29|CD .|.\nMr.|NNP Vinken|NNP is|VBZ chairman|NN of|IN\n
Elsevier|NNP N.V.|NNP ,|, the|DT Dutch|NNP publishing|VBG\n
group|NN .|. Rudolph|NNP Agnew|NNP ,|, 55|CD years|NNS old|JJ\n
and|CC former|JJ chairman|NN of|IN Consolidated|NNP Gold|NNP\n
Fields|NNP PLC|NNP ,|, was|VBD named|VBN a|DT nonexecutive|JJ\n
director|NN of|IN this|DT British|JJ industrial|JJ conglomerate|NN\n
.|.\nA|DT form|NN of|IN asbestos|NN once|RB used|VBN to|TO make|VB\n
Kent|NNP cigarette|NN filters|NNS has|VBZ caused|VBN a|DT high|JJ\n
percentage|NN of|IN cancer|NN deaths|NNS among|IN a|DT group|NN\n
of|IN workers|NNS exposed|VBN to|TO it|PRP more|RBR than|IN\n
30|CD years|NNS ago|IN ,|, researchers|NNS reported|VBD .|.`

func getSentences() []gonlp.WordsTags {

	var sentences []gonlp.WordsTags

	sent := strings.Split(TEXT, "\\n")

	for _, s := range sent {
		pairs := strings.Split(s, " ")
		wt := gonlp.WordsTags{}
		for _, w := range pairs {
			pair := strings.Split(w, "|")
			wt.Words = append(wt.Words, pair[0])
			wt.Tags = append(wt.Tags, pair[1])
		}
		sentences = append(sentences, wt)
	}
	return sentences
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	tokenizer := gonlp.NewSplitTokenizer(" ")

	//posTagger, err := gonlp.NewPerceptronTagger(tokenizer, false)
	posTagger, err := gonlp.NewPerceptronTagger(tokenizer, true)

	// posTagger.Train(getSentences(), "avp_model.gob", 5)

	if err != nil {
		fmt.Fprintln(os.Stderr, "posTagger:", err)
	}

	for scanner.Scan() {
		line := scanner.Text()

		tags, _ := posTagger.Tag(line)

		for _, tag := range tags {
			fmt.Println(tag)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
