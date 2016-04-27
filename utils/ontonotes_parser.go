package utils

import (
	"bufio"
	"fmt"
	"github.com/korobool/nlp4go/pos"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SentenceHandler func(string)
type FileHandler func(reader io.Reader) error

type OntonotesParser struct {
	OntonotesPath string
	reSpace       *regexp.Regexp
	reWordTag     *regexp.Regexp
}

func NewOntonotesParser(path string) (*OntonotesParser, error) {
	return &OntonotesParser{
		OntonotesPath: path,
		reSpace:       regexp.MustCompile(`\s+`),
		reWordTag:     regexp.MustCompile(`\(([^\s\(\)]+) ([^\s]+)\)`),
	}, nil
}
func (p *OntonotesParser) ParseToPath(outPath string) error {
	return p.parseToPath(p.OntonotesPath, outPath)
}

func (p *OntonotesParser) parseToPath(inPath, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	handleFn := func(reader io.Reader) error {
		return p.ParseFileToWriter(reader, writer)
	}
	return p.walkPath(inPath, handleFn)
}

func (p *OntonotesParser) ParseToWordsTags() ([]pos.WordsTags, error) {
	return p.parseToWordsTags(p.OntonotesPath)
}

func (p *OntonotesParser) parseToWordsTags(inPath string) ([]pos.WordsTags, error) {

	var sentences []pos.WordsTags

	handleFn := func(reader io.Reader) error {

		parsedSent, err := p.ParseFileToWordsTags(reader)
		if err == nil {
			sentences = append(sentences, parsedSent...)
		}
		return err
	}
	if err := p.walkPath(inPath, handleFn); err != nil {
		return nil, err
	}
	return sentences, nil
}

func (p *OntonotesParser) walkPath(ontoPath string, handler FileHandler) error {

	walkFn := func(path string, f os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".parse") {
			file, _ := os.Open(path)
			defer file.Close()

			err := handler(file)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := filepath.Walk(ontoPath, walkFn)
	if err != nil {
		return err
	}
	return nil
}

func (p *OntonotesParser) ParseFileToWordsTags(reader io.Reader) ([]pos.WordsTags, error) {
	wtsList := []pos.WordsTags{}

	handlerFn := func(sent string) {
		match := p.reWordTag.FindAllStringSubmatch(sent, -1)

		wts := pos.WordsTags{}
		for _, pair := range match {
			wts.Words = append(wts.Words, pair[2])
			wts.Tags = append(wts.Tags, pair[1])
		}
		wtsList = append(wtsList, wts)
	}

	if err := p.parseFile(reader, handlerFn); err != nil {
		return nil, err
	}
	return wtsList, nil
}

func (p *OntonotesParser) ParseFileToWriter(reader io.Reader, writer *bufio.Writer) error {

	handlerFn := func(sent string) {
		byteStr := []byte("")
		match := p.reWordTag.FindAllStringSubmatch(sent, -1)

		for _, pair := range match {
			str := fmt.Sprintf("(%s %s)", pair[1], pair[2])
			byteStr = append(byteStr, []byte(str)...)
		}
		fmt.Fprintln(writer, string(byteStr))
	}

	if err := p.parseFile(reader, handlerFn); err != nil {
		return err
	}
	return nil
}

func (p *OntonotesParser) parseFile(reader io.Reader, handler SentenceHandler) error {

	sent := []byte{}

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if scanner.Text() != "" {
			sent = append(sent, scanner.Bytes()...)
			continue
		}
		if len(sent) > 0 {
			s := p.reSpace.ReplaceAllLiteralString(string(sent), " ")
			handler(s)
		}
		sent = []byte{}
	}
	return scanner.Err()
}
