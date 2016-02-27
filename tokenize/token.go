package tokenize

import "fmt"

type Token struct {
	Text          []rune `json:"text"`
	Pos           int    `json:"pos"`
	IsQuoteStart  bool   `json:"is_quote_start"`
	IsQuoteEnd    bool   `json:"is_quote_end"`
	IsEllipsis    bool   `json:"is_ellipsis"`
	HasApostrophe bool   `json:"has_apostrophe"`
}

func NewToken(str []rune, posStart, length int) *Token {
	return &Token{
		Text: str[posStart:length],
		Pos:  posStart,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("<text=%s start=%d end=%d q_s=%v q_e=%v e=%v a=%v>",
		string(t.Text),
		t.Pos,
		t.PosEnd(),
		t.IsQuoteStart,
		t.IsQuoteEnd,
		t.IsEllipsis,
		t.HasApostrophe,
	)
}

func (t *Token) Len() int {
	return len(t.Text)
}

func (t *Token) PosEnd() int {
	return t.Pos + len(t.Text)
}

func (t *Token) Equals(compare *Token) bool {
	if string(t.Text) != string(compare.Text) {
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
