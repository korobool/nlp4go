package tokenize

import "regexp"

type LangContractions interface {
	Expand(*Token) ([]*Token, bool)
}

type EnglishContractions struct {
	resApostr  []*regexp.Regexp
	resGeneric []*regexp.Regexp
}

func NewEnglishContractions() *EnglishContractions {
	return &EnglishContractions{
		resApostr: []*regexp.Regexp{
			regexp.MustCompile(`(?i)^[^' ]+('s|'m|'d|'ll|'re|'ve|n't|')$`),
			regexp.MustCompile(`(?i)^d('ye)$`),
			regexp.MustCompile(`(?i)^mor('n)$`),
			regexp.MustCompile(`(?i)^'t(is)$`),
			regexp.MustCompile(`(?i)^'t(was)$`),
		},
		resGeneric: []*regexp.Regexp{
			regexp.MustCompile(`(?i)^can(not)$`),
			regexp.MustCompile(`(?i)^got(ta)$`),
			regexp.MustCompile(`(?i)^(?:gim|lem)(me)$`),
			regexp.MustCompile(`(?i)^(?:gon|wan)(na)$`),
		},
	}
}

func (c *EnglishContractions) Expand(token *Token) ([]*Token, bool) {

	if token.HasApostrophe {
		for _, re := range c.resApostr {
			tokens, ok := c.splitToken(re, token)
			if ok {
				return tokens, true
			}
		}
	}
	for _, re := range c.resGeneric {
		tokens, ok := c.splitToken(re, token)
		if ok {
			return tokens, true
		}
	}
	return nil, false
}

func (c *EnglishContractions) splitToken(re *regexp.Regexp, token *Token) ([]*Token, bool) {

	word := token.Word
	match := re.FindStringSubmatchIndex(word)

	if len(match) == 4 {
		boundL, _ := byteToRunePosition(word, match[2], match[3])
		tokens := []*Token{
			&Token{
				Pos: token.Pos,
			},
			&Token{
				Pos: token.Pos + boundL,
			},
		}
		tokens[0].SetText(token.Runes[:boundL])
		tokens[1].SetText(token.Runes[boundL:])

		return tokens, true
	}
	return nil, false
}
