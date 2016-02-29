package tokenize

type SplitTokenizer struct {
	delimiter []rune
}

func NewSplitTokenizer(delimiter string) *SplitTokenizer {
	return &SplitTokenizer{
		delimiter: []rune(delimiter),
	}
}

func (t *SplitTokenizer) Tokenize(str string) []*Token {

	s := []rune(str)
	sep := t.delimiter
	tokens := make([]*Token, 0, 5)

	if len(sep) == 0 {
		return tokens
	}

	c := sep[0]
	start := 0

	for i, n := 0, 0; i+len(sep) <= len(s); i++ {
		if s[i] == c && (len(sep) == 1 || isEqualRunes(s[i:i+len(sep)], sep)) {
			if i-start > 0 {
				tokens = append(tokens, NewToken(s, start, len(s[start:i])))
			}
			n++
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	if start > 0 {
		if len(s)-start > 0 {
			tokens = append(tokens, NewToken(s, start, len(s[start:])))
		}
	}
	if len(tokens) == 0 {
		tokens = append(tokens, NewToken(s, 0, len(s)))
	}
	return tokens
}
