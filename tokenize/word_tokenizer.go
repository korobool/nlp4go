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

func (t *TBWordTokenizer) Tokenize(s []rune) []*SentenceToken {

	var preparedToken *SentenceToken
	var isTokenPrepared bool

	tokens := make([]*SentenceToken, 0, 50)

	commitPrepared := func(posEnd int) {
		if isTokenPrepared {
			preparedToken.PosEnd = posEnd
			preparedToken.Text = s[preparedToken.PosStart:posEnd]
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

			// increase iterator counter because token length could be over than one char
			if token.PosEnd-pos > 1 {
				pos = token.PosEnd
			}
			continue SCAN
		}

		if pos == len(s)-1 {
			commitPrepared(len(s))
			continue
		}

		if !isTokenPrepared {
			preparedToken = &SentenceToken{PosStart: pos}
			isTokenPrepared = true
		}
		// Set HasApostrophe property to find token candiadtes with contractions easiely
		if current == '\'' {
			preparedToken.HasApostrophe = true
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

func (t *TBWordTokenizer) expandContractions(tokens []*SentenceToken) ([]*SentenceToken, bool) {
	var modified bool

	expandedTokens := make([]*SentenceToken, 0, len(tokens))

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

func (t *TBWordTokenizer) normalize(tokens []*SentenceToken) bool {

	var modified bool

	for _, token := range tokens {
		if token.IsQuoteStart {
			// replace starting quote with ``
			token.Text = []rune{'`', '`'}
			modified = true

		} else if token.IsQuoteEnd {
			// replace starting ending quote with ''
			token.Text = []rune{'\'', '\''}
			modified = true
		}
	}
	return modified
}
