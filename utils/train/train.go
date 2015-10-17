package main

import (
	"bufio"
	"fmt"
	//"github.com/davecheney/profile"
	"github.com/korobool/nlp4go"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()

	sentences := []gonlp.WordsTags{}

	walkFn := func(path string, f os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".parse") {
			file, _ := os.Open(path)
			defer file.Close()

			wts, err := parseFile(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			sentences = append(sentences, wts...)
		}
		return nil
	}

	filepath.Walk("/home/oleksandr/ontonotes", walkFn)

	posTagger, _ := gonlp.NewPerceptronTagger(nil, false, "")
	posTagger.Train(sentences, "avp_model.gob", 5)

}

func parseFile(reader io.Reader) ([]gonlp.WordsTags, error) {

	wts := []gonlp.WordsTags{}
	sent := []byte{}
	re := regexp.MustCompile("\\s+")

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if scanner.Text() != "" {
			sent = append(sent, scanner.Bytes()...)
			continue
		}
		if len(sent) > 0 {
			s := re.ReplaceAllLiteralString(string(sent), " ")
			wts = append(wts, makeWordsTags(s))
		}
		sent = []byte{}
	}
	err := scanner.Err()

	return wts, err
}

func makeWordsTags(tree string) gonlp.WordsTags {

	wt := gonlp.WordsTags{}
	re := regexp.MustCompile("\\(([^\\s\\(\\)]+) ([^\\s\\(\\)]+)\\)")

	match := re.FindAllStringSubmatch(tree, -1)
	for _, pair := range match {
		wt.Words = append(wt.Words, pair[2])
		wt.Tags = append(wt.Tags, pair[1])
	}

	return wt
}
