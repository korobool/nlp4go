package gonlp

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type PosTagger interface {
	Tag(string) []Tag
}

type Tag struct {
	Token Token
	Pos   string
}

type WordsTags struct {
	Words []string
	Tags  []string
}

type FeatureClass struct {
	f string
	c string
}

type BaseTagger struct {
	tokenizer Tokenizer
}

type AveragedPerceptron struct {
	Classes map[string]struct{}
	i       int
	totals  map[FeatureClass]float64
	tstamps map[FeatureClass]int
	weights map[string]map[string]float64
}

type PerceptronTagger struct {
	Model   *AveragedPerceptron
	TagMap  map[string]string
	classes map[string]struct{}
	BaseTagger
}

var (
	START          = []string{"-START-", "-START2-"}
	END            = []string{"-END-", "-END2-"}
	MODEL_GOB_PATH = "avp_model.gob"
)

func NewAveragedPerceptron() *AveragedPerceptron {
	var ap AveragedPerceptron

	ap.totals = make(map[FeatureClass]float64)
	ap.tstamps = make(map[FeatureClass]int)
	ap.weights = make(map[string]map[string]float64)

	return &ap
}

func (ap *AveragedPerceptron) Predict(features map[string]int) string {

	scores := map[string]float64{}

	for feat, value := range features {
		if _, ok := ap.weights[feat]; !ok || value == 0 {
			continue
		}
		weights := ap.weights[feat]

		for label, weight := range weights {
			scores[label] += float64(value) * weight
		}
	}

	return MaxScore(scores)
}

func (ap *AveragedPerceptron) Update(truth string, guess string, features map[string]int) {

	updFeat := func(c string, f string, w float64, v float64) {
		param := FeatureClass{f, c}

		ap.totals[param] += float64(ap.i-ap.tstamps[param]) * w
		ap.tstamps[param] = ap.i
		ap.weights[f][c] = w + v
	}

	ap.i += 1
	if truth == guess {
		return
	}
	for f := range features {
		weights, ok := ap.weights[f]
		if !ok {
			ap.weights[f] = make(map[string]float64)
			weights, _ = ap.weights[f]
		}
		wt, ok := weights[truth]
		if !ok {
			wt = 0.0
		}
		wg, ok := weights[guess]
		if !ok {
			wg = 0.0
		}

		updFeat(truth, f, wt, 1.0)
		updFeat(guess, f, wg, -1.0)
	}
	return
}

func (ap *AveragedPerceptron) AverageWeights() {

	for feat, weights := range ap.weights {

		newFeatWeights := map[string]float64{}
		for class, weight := range weights {
			param := FeatureClass{feat, class}

			total := ap.totals[param]
			total += float64(ap.i-ap.tstamps[param]) * weight

			averaged := Round(total/float64(ap.i), 0.5, 3)
			if averaged != 0 {
				newFeatWeights[class] = averaged
			}
		}
		ap.weights[feat] = newFeatWeights
	}
	return
}

func NewPerceptronTagger(tokenizer Tokenizer, load bool, path string) (*PerceptronTagger, error) {

	tagger := PerceptronTagger{}
	tagger.Model = NewAveragedPerceptron()
	tagger.classes = make(map[string]struct{})
	tagger.TagMap = make(map[string]string)
	tagger.tokenizer = tokenizer

	if load {
		if path == "" {
			path = MODEL_GOB_PATH
		}
		err := tagger.loadGob(path)
		if err != nil {
			return nil, err
		}
	}

	return &tagger, nil
}

func (t *PerceptronTagger) Tag(sentence string) ([]Tag, error) {
	var tags []Tag

	tokens := t.tokenizer.Tokenize(sentence)

	prev, prev2 := START[0], START[1]
	context := []string{}
	context = append(context, START...)
	for _, token := range tokens {
		context = append(context, normalize(token.Word))
	}
	context = append(context, END...)

	for i, token := range tokens {
		tag, ok := t.TagMap[token.Word]
		if !ok {
			features := getFeatures(i, token.Word, context, prev, prev2)
			tag = t.Model.Predict(features)
		}
		tags = append(tags, Tag{Token: token, Pos: tag})
		prev2 = prev
		prev = tag
	}

	return tags, nil
}

func (t *PerceptronTagger) Train(sentences []WordsTags, savePath string, iter int) {

	t.makeTagMap(&sentences)
	t.Model.Classes = t.classes

	prev, prev2 := START[0], START[1]

	fmt.Println("sentences len: ", len(sentences))
	for it := 0; it < iter; it += 1 {
		cnt := 1
		for _, wordsTags := range sentences {
			context := []string{}
			context = append(context, START...)

			for _, word := range wordsTags.Words {
				context = append(context, normalize(word))
			}
			context = append(context, END...)

			for i, word := range wordsTags.Words {

				guess, ok := t.TagMap[word]
				if !ok {
					features := getFeatures(i, word, context, prev, prev2)
					guess = t.Model.Predict(features)

					t.Model.Update(wordsTags.Tags[i], guess, features)
				}
				prev2 = prev
				prev = guess
			}
			fmt.Printf("\rsent # %d", cnt)
			cnt += 1
		}
		if (it + 1) < iter {
			for i, v := range rand.Perm(len(sentences)) {
				sentences[i], sentences[v] = sentences[v], sentences[i]
			}
		}
	}
	t.Model.AverageWeights()

	if savePath != "" {
		err := t.saveGob(savePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
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
			t.classes[pair[1]] = struct{}{}
		}
	}

	freqThresh := 20
	ambiguityThresh := 0.97

	for word, tagFreqs := range counts {

		var maxTag string
		var maxFreq, sum int

		for tag, freq := range tagFreqs {

			if freq > maxFreq {
				maxTag, maxFreq = tag, freq
			}
			sum += freq
		}

		if sum >= freqThresh && float64(maxFreq)/float64(sum) >= ambiguityThresh {
			t.TagMap[word] = maxTag
		}
	}
}

func (t *PerceptronTagger) loadGob(filename string) error {
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

	t.Model.weights = dump["weights"].(map[string]map[string]float64)
	t.TagMap = dump["tagmap"].(map[string]string)
	t.classes = dump["classes"].(map[string]struct{})
	t.Model.Classes = t.classes

	if err != nil {
		return err
	}
	return nil
}

func (t *PerceptronTagger) saveGob(filename string) error {

	gob.Register(map[string]map[string]float64{})
	gob.Register(map[string]string{})
	gob.Register(map[string]struct{}{})

	data := make(map[string]interface{})

	data["weights"] = t.Model.weights
	data["tagmap"] = t.TagMap
	data["classes"] = t.classes

	gobFile, err := os.Create(filename)
	defer gobFile.Close()

	if err != nil {
		return err
	}
	enc := gob.NewEncoder(gobFile)
	err = enc.Encode(data)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func normalize(word string) string {
	isDigit := func(s string) bool {
		for _, char := range s {
			if !(char >= '0' && char <= '9') {
				return false
			}
		}
		return true
	}
	switch {
	case strings.ContainsAny(word, "-") && !strings.HasPrefix(word, "-"):
		return "!HYPHEN"
	case len(word) == 4 && isDigit(word):
		return "!YEAR"
	case strings.IndexAny(word, "0123456789") == 0:
		return "!DIGITS"
	default:
		return strings.ToLower(word)
	}
}
func getFeatures(i int, word string, context []string, prev string, prev2 string) map[string]int {

	features := map[string]int{}

	add := func(name string, slice ...string) {
		str := []string{name}
		str = append(str, slice...)
		features[strings.Join(str, " ")] += 1
	}

	suffix := func(str string, sz int) string {
		rstr := []rune(str)
		if len(rstr) < sz {
			return str
		}
		return string(rstr[len(rstr)-sz:])
	}

	prefix := func(str string, sz int) string {
		rstr := []rune(str)
		if len(rstr) < sz {
			return str
		}
		return string(rstr[:len(rstr)-sz])
	}

	i += len(START)

	add("bias")
	add("i suffix", suffix(word, 3))
	add("i pref1", prefix(word, 1))
	add("i-1 tag", prev)
	add("i-2 tag", prev2)
	add("i tag+i-2 tag", prev, prev2)
	add("i word", context[i])
	add("i-1 tag+i word", prev, context[i])
	add("i-1 word", context[i-1])

	return features
}
