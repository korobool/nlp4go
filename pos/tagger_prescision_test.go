// +build prescision

package nlp4go

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/korobool/nlp4go/pos"
	"github.com/korobool/nlp4go/tokenize"
	"github.com/korobool/nlp4go/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	TmpPath string
)

const (
	parsedFileName   string = "parsed.corpus"
	validateFileName string = "validate.corpus"
	trainFileName    string = "train.corpus"
	MODEL_FILE_NAME  string = "test-model.gob"
)

func TestMain(m *testing.M) {

	atexit := func(exitCode int, exitMsg string) {
		if exitMsg != "" {
			log.Print(exitMsg)
		}
		os.RemoveAll(TmpPath)
		os.Exit(exitCode)
	}

	corpusPath := os.Getenv("ONTONOTES_PATH")
	if corpusPath == "" {
		atexit(1, "Enviroment variable ONTONOTES_PATH not set")
	}
	tmpPath, err := ioutil.TempDir(os.TempDir(), "tests-")
	if err != nil {
		atexit(1, fmt.Sprintf("Faield to create temporary directory: %v", err))
	}
	TmpPath = tmpPath
	if _, err := os.Stat(corpusPath); os.IsNotExist(err) {
		atexit(1, fmt.Sprintf("No such file or directory: %s", corpusPath))
	}
	parsedPath := filepath.Join(TmpPath, parsedFileName)
	if err := parseCorpus(corpusPath, parsedPath); err != nil {
		atexit(1, fmt.Sprintf("Failed to parse corpus: %v", err))
	}
	err = splitCorpus(TmpPath, parsedFileName, validateFileName, trainFileName)
	if err != nil {
		atexit(1, fmt.Sprintf("Failed to split file: %v", err))
	}
	if err := trainModel(TmpPath, trainFileName); err != nil {
		atexit(1, fmt.Sprintf("Faield to train model: %v", err))
	}

	atexit(m.Run(), "")
}

func TestPoSTaggerQuality(t *testing.T) {

	tokenizer := tokenize.NewSplitTokenizer(" ")

	cfgTagger := pos.TaggerConfig{
		Tokenizer: tokenizer,
		LoadModel: true,
		ModelPath: filepath.Join(TmpPath, MODEL_FILE_NAME),
	}
	posTagger, err := pos.NewPerceptronTagger(cfgTagger)
	if err != nil {
		t.Fatalf("Failed to create POS tagger: %v", err)
	}

	validateFilePath := filepath.Join(TmpPath, validateFileName)
	validateFile, err := os.Open(validateFilePath)
	if err != nil {
		t.Fatalf("Failed to open validation file: %v", err)
	}
	defer validateFile.Close()

	totalTags := 0
	guessedTags := 0

	scanner := bufio.NewScanner(validateFile)
	for scanner.Scan() {

		line := scanner.Text()
		wt := parseWordsTags(line)

		tokens, err := posTagger.Tag(strings.Join(wt.Words, " "))
		if err != nil {
			t.Fatalf("Error while tagging: %v", err)
		}
		if len(tokens) != len(wt.Tags) {
			t.Logf("Actual: %v", tokens)
			t.Logf("Expected: %v", wt.Tags)
			t.Error("Actual tokens count dosen't equal expected")
		}

		totalTags += len(wt.Tags)
		for index, knownTag := range wt.Tags {
			if knownTag == tokens[index].PosTag {
				guessedTags += 1
			}
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Error while reading validation file: %v", err)
	}
	t.Logf("Total: %d, guessed: %d, rate: %.2f", totalTags, guessedTags, float64(100*guessedTags)/float64(totalTags))

}

func trainModel(tmpPath, trainFileName string) error {

	trainFilePath := filepath.Join(tmpPath, trainFileName)

	trainSentences, err := parseWordsTagsFile(trainFilePath)
	if err != nil {
		return fmt.Errorf("Failed to parse train file: %v", err)
	}

	tokenizer := tokenize.NewTBWordTokenizer(true, true, nil)

	cfgTagger := pos.TaggerConfig{
		Tokenizer: tokenizer,
		LoadModel: false,
		ModelPath: filepath.Join(tmpPath, MODEL_FILE_NAME),
	}
	posTagger, err := pos.NewPerceptronTagger(cfgTagger)
	if err != nil {
		return fmt.Errorf("Failed to create POS tagger: %v", err)
	}
	posTagger.Train(trainSentences, 5, nil)
	posTagger.SaveModel()

	return nil
}

func parseWordsTagsFile(filePath string) ([]pos.WordsTags, error) {
	wts := []pos.WordsTags{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		wts = append(wts, parseWordsTags(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return wts, nil
}

func parseWordsTags(str string) pos.WordsTags {

	wt := pos.WordsTags{}
	re := regexp.MustCompile(`\(([^\s\(\)]+) ([^\s\(\)]+)\)`)

	match := re.FindAllStringSubmatch(str, -1)
	for _, pair := range match {
		wt.Words = append(wt.Words, pair[2])
		wt.Tags = append(wt.Tags, pair[1])
	}
	return wt
}

func parseCorpus(corpusPath, outputPath string) error {
	parser, err := utils.NewOntonotesParser(corpusPath)
	if err != nil {
		return err
	}
	if err := parser.ParseToPath(outputPath); err != nil {
		return err
	}
	return nil
}

func splitCorpus(tmpPath, parsedFileName, validateFileName, trainFileName string) error {

	parsedPath := filepath.Join(tmpPath, parsedFileName)
	validatePath := filepath.Join(tmpPath, validateFileName)
	trainPath := filepath.Join(tmpPath, trainFileName)

	parsedLines, err := countLines(parsedPath)
	if err != nil {
		return err
	}

	parsedFile, err := os.Open(parsedPath)
	if err != nil {
		return err
	}
	defer parsedFile.Close()

	validateFile, err := os.Create(validatePath)
	if err != nil {
		return err
	}
	defer validateFile.Close()

	trainFile, err := os.Create(trainPath)
	if err != nil {
		return err
	}
	defer trainFile.Close()

	scanner := bufio.NewScanner(parsedFile)
	//scanner.Split(bufio.ScanLines)

	validateWr := bufio.NewWriter(validateFile)
	defer validateWr.Flush()

	trainWr := bufio.NewWriter(trainFile)
	defer trainWr.Flush()

	for line := 0; scanner.Scan(); line++ {
		if line <= parsedLines/10 {
			fmt.Fprintln(validateWr, scanner.Text())
		} else {
			fmt.Fprintln(trainWr, scanner.Text())
		}
	}
	return scanner.Err()
}

func countLines(path string) (int, error) {

	r, err := os.Open(path)
	if err != nil {
		return 0, nil
	}
	defer r.Close()

	buf := make([]byte, 32768)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}
		count += bytes.Count(buf[:c], lineSep)
		if err == io.EOF {
			break
		}
	}
	return count, nil
}
