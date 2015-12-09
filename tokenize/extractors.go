package tokenize

import (
	"golang.org/x/text/unicode/rangetable"
	"unicode"
)

type TokenExtractor func([]rune, int) (*SentenceToken, bool)

var (
	startQuotesTbl = rangetable.New(' ', '(', '[', '{', '<')
	endPeriodTbl   = rangetable.New(']', '}', '}', '>', '\'', '"')
	standaloneTbl  = rangetable.New('?', '!', ';', '@', '#', '$', '%', '&', '(', ')', '[', ']', '{', '}', '<', '>')
)

func extractTokenQuote(s []rune, pos int) (*SentenceToken, bool) {

	if s[pos] != '"' {
		return nil, false
	}

	token := NewSentenceToken(s, pos, pos+1)

	if pos == 0 {
		token.IsQuoteStart = true
	} else if unicode.Is(startQuotesTbl, s[pos-1]) {
		token.IsQuoteStart = true
	} else {
		token.IsQuoteEnd = true
	}

	return token, true
}

func extractTokenPeriod(s []rune, pos int) (*SentenceToken, bool) {

	if s[pos] != '.' {
		return nil, false
	}

	var token *SentenceToken

	if pos != 0 && s[pos-1] == '.' {
		return nil, false
	}

	if pos == len(s)-1 {
		token = NewSentenceToken(s, pos, pos+1)

	} else if len(s) > pos+2 && s[pos+1] == '.' && s[pos+2] == '.' {
		// NOTE for ellipsis we should increment position
		token = NewSentenceToken(s, pos, pos+3)
		token.IsEllipsis = true

	} else if unicode.Is(endPeriodTbl, s[pos+1]) || unicode.IsSpace(s[pos+1]) {
		// NOTE checking all chars till the end whether it
		// closing parens/brackets, etc. BUT not more than 5 chars examined
		for maxIt, i := 5, pos+1; i < len(s) || i <= maxIt; i++ {
			if !unicode.Is(endPeriodTbl, s[i]) && !unicode.IsSpace(s[i]) {
				return nil, false
			}
		}
		token = NewSentenceToken(s, pos, pos+1)

	} else {
		return nil, false
	}

	return token, true
}

func extractTokenApostrophe(s []rune, pos int) (*SentenceToken, bool) {
	if s[pos] != '\'' {
		return nil, false
	}
	if pos == 0 || pos == len(s)-1 {
		return nil, false

	} else if s[pos-1] != '\'' && unicode.IsSpace(s[pos+1]) {
		return NewSentenceToken(s, pos, pos+1), true
	}

	return nil, false
}

func extractTokenColon(s []rune, pos int) (*SentenceToken, bool) {
	if s[pos] != ':' {
		return nil, false
	}

	if pos != len(s)-1 && unicode.IsDigit(s[pos+1]) {
		return nil, false
	}
	return NewSentenceToken(s, pos, pos+1), true
}

func extractTokenComma(s []rune, pos int) (*SentenceToken, bool) {
	if s[pos] != ',' {
		return nil, false
	}

	if pos != len(s)-1 && unicode.IsDigit(s[pos+1]) {
		return nil, false
	}
	return NewSentenceToken(s, pos, pos+1), true
}

func extractTokenHyphen(s []rune, pos int) (*SentenceToken, bool) {
	if s[pos] != '-' {
		return nil, false
	}
	if pos != len(s)-1 && s[pos+1] == '-' {
		// NOTE for double-hyphen we should increment position
		return NewSentenceToken(s, pos, pos+2), true
	}
	return nil, false
}

func extractTokenSymbol(s []rune, pos int) (*SentenceToken, bool) {
	if !unicode.Is(standaloneTbl, s[pos]) {
		return nil, false
	}
	return NewSentenceToken(s, pos, pos+1), true
}
