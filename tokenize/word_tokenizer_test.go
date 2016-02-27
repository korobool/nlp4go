package tokenize

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

type RefSentence struct {
	Sentence string   `json:"sentence"`
	Tokens   []*Token `json:"tokens"`
}

func loadRefSentences(fileName string) ([]RefSentence, error) {
	var refList []RefSentence

	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &refList)
	if err != nil {
		return nil, err
	}
	return refList, nil
}

func TestEnglishTokensCount(t *testing.T) {
	refList, err := loadRefSentences("../test_data/sentences.en.json")
	if err != nil {
		t.Fatal(err)
	}

	tokenizer := NewTBWordTokenizer(true, true, nil)

	for line, refSent := range refList {

		tokens := tokenizer.Tokenize([]rune(refSent.Sentence))

		if len(tokens) != len(refSent.Tokens) {
			t.Logf("Actual  : len=%d tokens=%v", len(tokens), tokens)
			t.Log("--------")
			t.Logf("Expected: len=%d tokens=%v", len(refSent.Tokens), refSent.Tokens)
			t.Fatalf("Line #%d: actual token count != expected", line)
		}
		for i, token := range tokens {
			if !refSent.Tokens[i].Equals(token) {
				t.Logf("Actual  : %v", token)
				t.Log("--------")
				t.Logf("Expected: %v", refSent.Tokens[i])
				t.Logf("Line #%d: actual token #%d != expected", line, i)
				t.Fatalf("Tokens  : %v", tokens)
			}
		}
	}

}
