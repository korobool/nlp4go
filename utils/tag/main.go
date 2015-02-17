package main

import (
	"bufio"
	"fmt"
	"github.com/korobool/go-nlp-tools"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	tokenizer := gonlp.NewSplitTokenizer(" ")

	posTagger, err := gonlp.NewPerceptronTagger(tokenizer, true)

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
