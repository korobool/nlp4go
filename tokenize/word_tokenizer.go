package tokenize

import "unicode"

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
		preparedToken.PosEnd = posEnd
		preparedToken.Text = s[preparedToken.PosStart:posEnd]
		isTokenPrepared = false
		tokens = append(tokens, preparedToken)
	}

SCAN:
	for pos := 0; pos < len(s); pos++ {

		current := s[pos]

		if unicode.IsSpace(current) {
			if isTokenPrepared {
				commitPrepared(pos)
			}
			continue
		}

		for _, extractFn := range t.extractors {
			token, ok := extractFn(s, pos)
			if !ok {
				continue
			}
			if isTokenPrepared {
				commitPrepared(pos)
			}
			tokens = append(tokens, token)
			// increase iterator counter because token length could be over than one char
			if token.PosEnd-pos > 1 {
				pos = token.PosEnd
			}
			continue SCAN
		}

		if isTokenPrepared && pos == len(s)-1 {
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

	if t.ExpandContrations {
		if expandedTokens, ok := t.expandContractions(tokens); ok {
			return expandedTokens
		}
	}

	return tokens
}

func (t *TBWordTokenizer) expandContractions(tokens []*SentenceToken) ([]*SentenceToken, bool) {
	var tokensModified bool

	expandedTokens := make([]*SentenceToken, 0, len(tokens))

	for i, tok := range tokens {
		if tokenList, ok := t.LangContractions.Expand(tok); ok {
			if i > 0 && !tokensModified {
				expandedTokens = append(expandedTokens, tokens[:i-1]...)
			}
			tokensModified = true
			expandedTokens = append(expandedTokens, tokenList...)

		} else if tokensModified {
			expandedTokens = append(expandedTokens, tok)
		}
	}
	if tokensModified {
		return expandedTokens, true
	}

	return tokens, false
}
