package pos

import (
	"encoding/gob"
	"github.com/korobool/nlp4go/ml"
	"github.com/korobool/nlp4go/tokenize"
	"math/rand"
	"os"
	"strings"
	"unicode"
)

// callback func(current, total) for Train progress
type Callback func(int, int)

type PerceptronTagger struct {
	tokenizer          tokenize.Tokenizer
	FrequencyThreshold int
	AmbiguityThreshold float64
	START_TOK          []string
	END_TOK            []string
	ModelPath          string
	Model              *ml.AveragedPerceptron
	TagMap             map[string]string
	Classes            map[string]struct{}
}

func NewPerceptronTagger(config TaggerConfig) (*PerceptronTagger, error) {

	tagger := PerceptronTagger{
		tokenizer:          config.Tokenizer,
		FrequencyThreshold: 20,
		AmbiguityThreshold: 0.97,
		START_TOK:          []string{"-START-", "-START2-"},
		END_TOK:            []string{"-END-", "-END2-"},
		ModelPath:          "avp_tagger_model.gob",
		Model:              ml.NewAveragedPerceptron(),
		TagMap:             make(map[string]string),
		Classes:            make(map[string]struct{}),
	}

	if config.ModelPath != "" {
		tagger.ModelPath = config.ModelPath
	}
	if config.LoadModel {
		err := tagger.LoadModel(tagger.ModelPath)
		if err != nil {
			return nil, err
		}
	}
	return &tagger, nil
}

func (t *PerceptronTagger) Tag(sentence string) ([]*tokenize.Token, error) {

	tokens := t.tokenizer.Tokenize(sentence)

	prev, prev2 := t.START_TOK[0], t.START_TOK[1]

	context := make([]string, 0, len(tokens)+4)
	context = append(context, t.START_TOK...)

	for _, token := range tokens {
		context = append(context, t.normalize(token.Word, token.Runes))
	}
	context = append(context, t.END_TOK...)

	for i, token := range tokens {
		tag, ok := t.TagMap[token.Word]
		if !ok {
			features := t.getFeatures(i, token.Word, context, prev, prev2)
			tag = t.Model.Predict(features)
		}
		tokens[i].PosTag = tag
		prev2 = prev
		prev = tag
	}
	return tokens, nil
}

func (t *PerceptronTagger) Train(sentences []WordsTags, rounds int, progressFn Callback) {

	t.makeTagMap(&sentences)

	prev, prev2 := t.START_TOK[0], t.START_TOK[1]

	for it := 0; it < rounds; it += 1 {
		cnt := 1
		for _, wordsTags := range sentences {

			context := make([]string, 0, len(wordsTags.Words)+4)
			context = append(context, t.START_TOK...)

			for _, word := range wordsTags.Words {
				context = append(context, t.normalize(word, []rune(word)))
			}
			context = append(context, t.END_TOK...)

			for i, word := range wordsTags.Words {
				guess, ok := t.TagMap[word]
				if !ok {
					features := t.getFeatures(i, word, context, prev, prev2)
					guess = t.Model.Predict(features)

					t.Model.Update(wordsTags.Tags[i], guess, features)
				}
				prev2 = prev
				prev = guess
			}
			if progressFn != nil {
				progressFn((it+1)*cnt, len(sentences)*rounds)
			}
			cnt += 1
		}
		if (it + 1) < rounds {
			for i, v := range rand.Perm(len(sentences)) {
				sentences[i], sentences[v] = sentences[v], sentences[i]
			}
		}
	}
	t.Model.AverageWeights()
}

func (t *PerceptronTagger) makeTagMap(sentences *[]WordsTags) {

	counts := make(map[string]map[string]int)

	for _, sent := range *sentences {

		wordTag := make([][2]string, len(sent.Words), len(sent.Words))

		for i, word := range sent.Words {
			wordTag[i] = [2]string{word, sent.Tags[i]}
		}
		for _, pair := range wordTag {
			if _, ok := counts[pair[0]]; !ok {
				counts[pair[0]] = make(map[string]int)
			}
			counts[pair[0]][pair[1]] += 1
			t.Classes[pair[1]] = struct{}{}
		}
	}

	for word, tagFreqs := range counts {
		var maxTag string
		var maxFreq, sum int

		for tag, freq := range tagFreqs {
			if freq > maxFreq {
				maxTag, maxFreq = tag, freq
			}
			sum += freq
		}
		if sum >= t.FrequencyThreshold &&
			(float64(maxFreq)/float64(sum) >= t.AmbiguityThreshold) {
			t.TagMap[word] = maxTag
		}
	}
}

func (t *PerceptronTagger) normalize(word string, runes []rune) string {

	switch {
	case strings.ContainsAny(word, "-") && !strings.HasPrefix(word, "-"):
		return "!HYPHEN"
	case len(runes) == 4 && IsDigit(word):
		return "!YEAR"
	case unicode.IsDigit(runes[0]):
		return "!DIGITS"
	default:
		return strings.ToLower(word)
	}
}

func (t *PerceptronTagger) getFeatures(i int, word string, context []string, prev string, prev2 string) map[string]int {

	features := map[string]int{}

	add := func(name string, slice ...string) {
		str := make([]string, 0, len(slice)+1)
		str = append(str, name)
		str = append(str, slice...)
		features[strings.Join(str, " ")] += 1
	}

	i += len(t.START_TOK)

	add("bias")
	add("i suffix", StrSuffix(word, 3))
	add("i pref1", StrPrefix(word, 1))
	add("i-1 tag", prev)
	add("i-2 tag", prev2)
	add("i tag+i-2 tag", prev, prev2)
	add("i word", context[i])
	add("i-1 tag+i word", prev, context[i])
	add("i-1 word", context[i-1])
	add("i-1 suffix", StrSuffix(context[i-1], 3))
	add("i-2 word", context[i-2])
	add("i+1 word", context[i+1])
	add("i+1 suffix", StrSuffix(context[i+1], 3))
	add("i+2 word", context[i+2])

	return features
}

func (t *PerceptronTagger) LoadModel(filename string) error {
	gob.Register(map[string]map[string]float64{})
	gob.Register(map[string]string{})
	gob.Register(map[string]struct{}{})

	dump := make(map[string]interface{})

	gobFile, err := os.Open(filename)
	defer gobFile.Close()
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(gobFile)
	err = dec.Decode(&dump)

	t.Model.Weights = dump["weights"].(map[string]map[string]float64)
	t.TagMap = dump["tagmap"].(map[string]string)
	t.Classes = dump["classes"].(map[string]struct{})

	return nil
}

func (t *PerceptronTagger) SaveModel() error {

	gob.Register(map[string]map[string]float64{})
	gob.Register(map[string]string{})
	gob.Register(map[string]struct{}{})

	data := make(map[string]interface{})

	data["weights"] = t.Model.Weights
	data["tagmap"] = t.TagMap
	data["classes"] = t.Classes

	gobFile, err := os.Create(t.ModelPath)
	defer gobFile.Close()
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(gobFile)
	if err := enc.Encode(data); err != nil {
		return err
	}

	return nil
}
