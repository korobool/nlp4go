// +build ignore

package main

import (
	"bufio"
	"fmt"
	"github.com/korobool/nlp4go"
	"os"
)

const (
	sentence1 string = "The dog saw a man in the park"
	sentence2 string = "The angry bear chased the frightened little squirrel"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	tokenizer := gonlp.NewSplitTokenizer(" ")
	posTagger, err := gonlp.NewPerceptronTagger(tokenizer, true, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "posTagger:", err)
	}

	syntaxParser := gonlp.NewSyntaxParser()

	for scanner.Scan() {
		var taggedTokens []gonlp.TaggedToken
		line := scanner.Text()
		tags, _ := posTagger.Tag(line)
		for _, tag := range tags {
			taggedTokens = append(
				taggedTokens,
				gonlp.TaggedToken{Token: tag.Token.Word, Tag: tag.Pos})
		}
		tree, _ := syntaxParser.Parse(taggedTokens)
		syntaxParser.PrettyPrint(tree, 0)
	}

	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
