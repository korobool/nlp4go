package gonlp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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
	parsedFileName   string = "pasred_corpus"
	validateFileName string = "validate_corpus"
	trainFileName    string = "train_corpus"
)

func TestMain(m *testing.M) {

	atexit := func(exitcode int) {
		//os.RemoveAll(TmpPath)
		os.Exit(exitcode)
	}

	corpusPath := os.Getenv("GONLP_CORPUS")
	if corpusPath == "" {
		fmt.Println("enviroment variable GONLP_CORPUS not set")
		atexit(1)
	}

	path, err := ioutil.TempDir(os.TempDir(), "tests-")
	if err != nil {
		fmt.Printf("faild to create temporary directory: %s \n", err)
		atexit(1)
	}
	TmpPath = path

	if _, err := os.Stat(corpusPath); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s \n", corpusPath)
		atexit(1)
	}

	parsedPath := filepath.Join(TmpPath, parsedFileName)
	if err := ParseOntonotes(corpusPath, parsedPath); err != nil {
		fmt.Printf("faild to parse corpus: %s \n", err)
		atexit(1)
	}

	err = splitCorpus()
	if err != nil {
		fmt.Printf("faild to split file: %s \n", err)
		atexit(1)
	}

	err = trainModel()
	if err != nil {
		fmt.Printf("faild to train model: %s \n", err)
		atexit(1)
	}

	atexit(m.Run())
}

func TestTagger(t *testing.T) {

	tokenizer := NewSplitTokenizer(" ")

	trainFilePath := filepath.Join(TmpPath, "avp_model.gob")
	posTagger, err := NewPerceptronTagger(tokenizer, true, trainFilePath)
	if err != nil {
		t.Errorf(err.Error())
	}

	validateFilePath := filepath.Join(TmpPath, validateFileName)
	validateFile, err := os.Open(validateFilePath)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer validateFile.Close()

	totalTags := 0
	guessedTags := 0

	scanner := bufio.NewScanner(validateFile)
	for scanner.Scan() {

		line := scanner.Text()
		wt := parseWordsTags(line)

		//fmt.Printf("WT: %q\n", strings.Join(wt.Words, " "))

		tags, err := posTagger.Tag(strings.Join(wt.Words, " "))
		if err != nil {
			t.Errorf(err.Error())
		}

		totalTags += len(wt.Tags)
		//fmt.Printf("[(%d) %q => (%d) %q]\n", len(wt.Tags), wt.Tags, len(tags), tags)
		for index, knownTag := range wt.Tags {

			if knownTag == tags[index].Pos {
				guessedTags += 1
			}
		}
	}
	if err := scanner.Err(); err != nil {
		t.Errorf(err.Error())
	}
	fmt.Printf("total: %d guessed: %d", totalTags, guessedTags)

}

func trainModel() error {

	posTagger, err := NewPerceptronTagger(nil, false, "")
	if err != nil {
		return err
	}

	sentences := []WordsTags{}

	trainFilePath := filepath.Join(TmpPath, trainFileName)
	trainFile, err := os.Open(trainFilePath)
	if err != nil {
		return err
	}
	defer trainFile.Close()

	scanner := bufio.NewScanner(trainFile)
	for scanner.Scan() {
		line := scanner.Text()
		sentences = append(sentences, parseWordsTags(line))
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	posTagger.Train(sentences, filepath.Join(TmpPath, "avp_model.gob"), 5)

	return nil
}

func parseWordsTags(str string) WordsTags {

	wt := WordsTags{}
	re := regexp.MustCompile(`\(([^\s\(\)]+) ([^\s\(\)]+)\)`)

	match := re.FindAllStringSubmatch(str, -1)
	for _, pair := range match {
		wt.Words = append(wt.Words, pair[2])
		wt.Tags = append(wt.Tags, pair[1])
	}
	return wt
}

func splitCorpus() error {

	parsedPath := filepath.Join(TmpPath, parsedFileName)
	validatePath := filepath.Join(TmpPath, validateFileName)
	trainPath := filepath.Join(TmpPath, trainFileName)

	parsedLines, err := lineCounter(parsedPath)
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

func lineCounter(path string) (int, error) {

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
