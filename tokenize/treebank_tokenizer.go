package tokenize

import (
	"unicode"
)

/*
	Mimics TreeBank word tokenizer without using mass of regexps
*/
type TBWordTokenizer struct {
	extractors        []TokenExtractor
	LangContractions  LangContractions
	ExpandContrations bool
	Normalize         bool
}

func NewTBWordTokenizer(normalize, checkContr bool, langContr LangContractions) *TBWordTokenizer {

	if langContr == nil && checkContr {
		langContr = NewEnglishContractions()
	}
	return &TBWordTokenizer{
		extractors: []TokenExtractor{
			extractTokenQuote,
			extractTokenPeriod,
			extractTokenApostrophe,
			extractTokenColon,
			extractTokenComma,
			extractTokenHyphen,
			extractTokenSymbol,
		},
		ExpandContrations: checkContr,
		LangContractions:  langContr,
		Normalize:         normalize,
	}
}

func (t *TBWordTokenizer) Tokenize(s string) []*Token {
	return t.TokenizeRune([]rune(s))
}

func (t *TBWordTokenizer) TokenizeRune(s []rune) []*Token {

	var preparedToken *Token
	var isTokenPrepared bool

	tokens := make([]*Token, 0, 50)

	commitPrepared := func(posEnd int) {
		if isTokenPrepared {
			preparedToken.SetText(s[preparedToken.Pos:posEnd])
			isTokenPrepared = false
			tokens = append(tokens, preparedToken)
		}
	}

SCAN:
	for pos := 0; pos < len(s); pos++ {

		current := s[pos]

		if unicode.IsSpace(current) {
			commitPrepared(pos)
			continue
		}

		for _, extractFn := range t.extractors {
			token, ok := extractFn(s, pos)
			if !ok {
				continue
			}
			commitPrepared(pos)
			tokens = append(tokens, token)

			// increase iterator counter because token length
			// could be over than one char
			if token.PosEnd()-pos > 1 {
				pos = token.PosEnd() - 1
			}
			continue SCAN
		}
		if !isTokenPrepared {
			preparedToken = &Token{Pos: pos}
			isTokenPrepared = true
		}

		// Set HasApostrophe property to find token candiadtes
		// with contractions easiely
		if current == '\'' {
			preparedToken.HasApostrophe = true
		}

		if pos == len(s)-1 {
			commitPrepared(len(s))
		}
	}

	if t.Normalize {
		t.normalize(tokens)
	}
	if t.ExpandContrations {
		if expandedTokens, ok := t.expandContractions(tokens); ok {
			return expandedTokens
		}
	}
	return tokens
}

func (t *TBWordTokenizer) expandContractions(tokens []*Token) ([]*Token, bool) {
	var modified bool

	expandedTokens := make([]*Token, 0, len(tokens))

	for i, tok := range tokens {
		if tokenList, ok := t.LangContractions.Expand(tok); ok {
			if i > 0 && !modified {
				expandedTokens = append(expandedTokens, tokens[:i]...)
			}
			modified = true
			expandedTokens = append(expandedTokens, tokenList...)

		} else if modified {
			expandedTokens = append(expandedTokens, tok)
		}
	}
	if modified {
		return expandedTokens, modified
	}
	return tokens, modified
}

func (t *TBWordTokenizer) normalize(tokens []*Token) bool {

	var modified bool

	for _, token := range tokens {
		if token.IsQuoteStart {
			// replace starting quote with ``
			token.SetText([]rune{'`', '`'})
			modified = true

		} else if token.IsQuoteEnd {
			// replace starting ending quote with ''
			token.SetText([]rune{'\'', '\''})
			modified = true
		}
	}
	return modified
}
