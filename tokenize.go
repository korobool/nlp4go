package gonlp

import (
	//"fmt"
	"regexp"
	"strings"
)

type Token struct {
	Word  string
	Start int
	End   int
}

type Tokenizer interface {
	Tokenize(string) []Token
}

type RegexpTokenizer struct {
	string
}

type SplitTokenizer struct {
	string
}

type TreeBankTokenizer struct {
	boundsRe *regexp.Regexp
	startQRe *regexp.Regexp
	res      []*regexp.Regexp
}

func NewSplitTokenizer(sep string) *SplitTokenizer {
	var t SplitTokenizer
	t.string = sep
	return &t
}

func NewRegexpTokenizer(pattern string) *RegexpTokenizer {
	var t RegexpTokenizer
	t.string = pattern
	return &t
}

func NewTreeBankTokenizer() *TreeBankTokenizer {

	var t TreeBankTokenizer

	t.boundsRe = regexp.MustCompile(`\S+`)
	t.startQRe = regexp.MustCompile(`(?:(?:^)|(?:[ (\[{<]))(")`)
	//t.endQRe = regexp.MustCompile(`(?:[^ (\[{<])(")`)

	t.res = []*regexp.Regexp{}

	t.res = append(t.res, regexp.MustCompile("(``)"))
	t.res = append(t.res, regexp.MustCompile(`([:,])[^\d]`))
	t.res = append(t.res, regexp.MustCompile(`(\.\.\.)`))
	t.res = append(t.res, regexp.MustCompile(`([;@#$%&?!])`))
	t.res = append(t.res, regexp.MustCompile(`[^\.](\.)[\]\)}>"']*\s*$`))
	t.res = append(t.res, regexp.MustCompile(`[^'](') `))
	t.res = append(t.res, regexp.MustCompile(`([\]\[\(\)\{\}\<\>])`))
	t.res = append(t.res, regexp.MustCompile(`(--)`))
	t.res = append(t.res, regexp.MustCompile(`(")`))
	// used \W instead of single space here
	t.res = append(t.res, regexp.MustCompile(`([^' ]+)('[sS]|'[mM]|'[dD]|'ll|'LL|'re|'RE|'ve|'VE|n't|N'T|')\W`))
	t.res = append(t.res, regexp.MustCompile(`(?i)\b(can)(not)\b`))
	t.res = append(t.res, regexp.MustCompile(`(?i)\b(d)('ye)\b`))
	t.res = append(t.res, regexp.MustCompile(`(?i)\b(gim)(me)\b`))
	t.res = append(t.res, regexp.MustCompile(`(?i)\b(gon)(na)\b`))
	t.res = append(t.res, regexp.MustCompile(`(?i)\b(got)(ta)\b`))
	t.res = append(t.res, regexp.MustCompile(`(?i)\b(lem)(me)\b`))
	t.res = append(t.res, regexp.MustCompile(`(?i)\b(mor)('n)\b`))
	t.res = append(t.res, regexp.MustCompile(`(?i)\b(wan)(na) `))
	t.res = append(t.res, regexp.MustCompile(`(?i) ('t)(is)\b`))
	t.res = append(t.res, regexp.MustCompile(`(?i) ('t)(was)\b`))
	t.res = append(t.res, regexp.MustCompile(`\S('')`))

	return &t
}

func Tokenize(text string, tokenizer Tokenizer) []Token {
	return tokenizer.Tokenize(text)
}

func (t *RegexpTokenizer) Tokenize(s string) []Token {

	tokens := []Token{}

	if len(t.string) > 0 && len(s) == 0 {
		return tokens
	}
	re := regexp.MustCompile(t.string)
	matches := re.FindAllStringIndex(s, -1)

	beg := 0
	end := 0
	for _, match := range matches {
		end = match[0]
		if match[1] != 0 {
			tokens = append(tokens, Token{Word: s[beg:end], Start: beg, End: end - 1})
		}
		beg = match[1]
	}
	if end != len(s) {
		tokens = append(tokens, Token{Word: s[beg:], Start: beg, End: len(s) - 1})
	}
	return tokens
}

func (t *SplitTokenizer) Tokenize(s string) []Token {

	sep := t.string
	tokens := []Token{}

	if sep == "" {
		//FIXME: error ?
		return tokens
	}
	n := strings.Count(s, sep) + 1
	c := sep[0]
	start := 0
	na := 0
	if n == 1 {
		tokens = append(tokens, Token{Word: s, Start: 0, End: len(s)})
		return tokens
	}
	for i := 0; i+len(sep) <= len(s) && na+1 < n; i++ {
		if s[i] == c && (len(sep) == 1 || s[i:i+len(sep)] == sep) {
			if i-start > 0 {
				tokens = append(tokens, Token{Word: s[start:i], Start: start, End: i})
			}
			na++
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	if start > 0 {
		if len(s)-start > 0 {
			tokens = append(tokens, Token{Word: s[start:], Start: start, End: len(s)})
		}
	}
	return tokens
}

func (t *TreeBankTokenizer) Tokenize(s string) []Token {

	tokens := []Token{}
	parts := [][]int{}
	startQuotes := [][]int{}

	pushToken := func(st int, en int, wo string) {

		if wo == `"` {
			for _, q := range startQuotes {
				if q[0] == st && q[1] == en {
					wo = "``"
				} else {
					wo = "''"
				}
			}
		}
		token := Token{Start: st, End: en, Word: wo}
		tokens = append(tokens, token)
	}

	for _, re := range t.res {
		parts = append(parts, getMatchedParts(s, re)...)
	}

	points := sortPoints(parts)
	bounds := t.boundsRe.FindAllStringSubmatchIndex(s, -1)

	startQuotes = getMatchedParts(s, t.startQRe)

	start := 0
	for _, bound := range bounds {
		prev := bound[0]
		for i := start; i < len(points) && points[i] <= bound[1]; i += 1 {
			if points[i] > prev {
				pushToken(prev, points[i], s[prev:points[i]])
			}
			prev = points[i]
			start += 1

		}
		if prev != bound[1] {
			pushToken(prev, bound[1], s[prev:bound[1]])
		}
	}
	return tokens
}

func getMatchedParts(s string, re *regexp.Regexp) [][]int {

	parts := [][]int{}

	matches := re.FindAllStringSubmatchIndex(s, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			for indx := 2; indx < len(match); indx += 2 {
				parts = append(parts, match[indx:indx+2])
			}
		}
	}
	return parts
}

func sortPoints(s [][]int) []int {

	points := []int{}

	for i := 1; i < len(s); i += 1 {
		for j := i; j > 0 && s[j][0] < s[j-1][0]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
	for _, v := range s {
		for i, _ := range v {
			if len(points) < 1 || (v[i] != points[len(points)-1] && v[i] > points[len(points)-1]) {
				points = append(points, v[i])
			}
		}
	}

	return points
}
