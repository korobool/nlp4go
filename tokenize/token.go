package tokenize

import "fmt"

type Token struct {
	Runes         []rune `json:"runes"`
	Word          string `json:"word"`
	Pos           int    `json:"pos"`
	PosTag        string `json:"pos_tag"`
	IsQuoteStart  bool   `json:"is_quote_start"`
	IsQuoteEnd    bool   `json:"is_quote_end"`
	IsEllipsis    bool   `json:"is_ellipsis"`
	HasApostrophe bool   `json:"has_apostrophe"`
}

func NewToken(str []rune, posStart, length int) *Token {
	return &Token{
		Runes: str[posStart : posStart+length],
		Word:  string(str[posStart : posStart+length]),
		Pos:   posStart,
	}
}

func (t *Token) SetText(text []rune) {
	t.Runes = text
	t.Word = string(text)
}

func (t *Token) String() string {
	return fmt.Sprintf("<word=%s pos_tag=%s start=%d end=%d q_s=%v q_e=%v e=%v a=%v>",
		t.Word,
		t.PosTag,
		t.Pos,
		t.PosEnd(),
		t.IsQuoteStart,
		t.IsQuoteEnd,
		t.IsEllipsis,
		t.HasApostrophe,
	)
}

func (t *Token) Len() int {
	return len(t.Runes)
}

func (t *Token) PosEnd() int {
	return t.Pos + len(t.Runes) - 1
}

func (t *Token) Equals(compare *Token) bool {
	if string(t.Runes) != string(compare.Runes) {
		return false
	}
	if t.Pos != compare.Pos {
		return false
	}
	if t.IsQuoteStart != compare.IsQuoteStart {
		return false
	}
	if t.IsEllipsis != compare.IsEllipsis {
		return false
	}
	if t.IsQuoteEnd != compare.IsQuoteEnd {
		return false
	}

	return true
}
