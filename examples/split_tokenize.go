// +build ignore

package main

import (
	"bufio"
	"fmt"
	"github.com/korobool/nlp4go/tokenize"
	"os"
)

func main() {

	tokenizer := tokenize.NewSplitTokenizer(" ")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		tokens := tokenizer.Tokenize(line)
		for _, token := range tokens {
			fmt.Println(token)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading standard input:", err)
	}
}
