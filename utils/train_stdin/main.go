package main

import (
	"fmt"
	//"github.com/davecheney/profile"
	"bufio"
	"github.com/korobool/go-nlp-tools"
	"os"
	"regexp"
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()

	sentences := []gonlp.WordsTags{}

	posTagger, err := gonlp.NewPerceptronTagger(nil, false)
	if err != nil {
		fmt.Fprintln(os.Stderr, "posTagger:", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		sentences = append(sentences, parseWordsTags(line))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	posTagger.Train(sentences, "avp_model.gob", 5)
}

func parseWordsTags(str string) gonlp.WordsTags {

	wt := gonlp.WordsTags{}
	re := regexp.MustCompile(`\(([^\s\(\)]+) ([^\s\(\)]+)\)`)

	match := re.FindAllStringSubmatch(str, -1)
	for _, pair := range match {
		wt.Words = append(wt.Words, pair[2])
		wt.Tags = append(wt.Tags, pair[1])
	}
	return wt
}
