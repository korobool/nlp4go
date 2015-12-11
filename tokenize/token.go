package tokenize

import "fmt"

type SentenceToken struct {
	Text          []rune `json:"text"`
	PosStart      int    `json:"pos_start"`
	PosEnd        int    `json:"pos_end"`
	IsQuoteStart  bool   `json:"is_quote_start"`
	IsQuoteEnd    bool   `json:"is_quote_end"`
	IsEllipsis    bool   `json:"is_ellipsis"`
	HasApostrophe bool   `json:"has_apostrophe"`
}

func NewSentenceToken(str []rune, posStart, posEnd int) *SentenceToken {
	return &SentenceToken{
		Text:     str[posStart:posEnd],
		PosStart: posStart,
		PosEnd:   posEnd,
	}
}

func (t *SentenceToken) String() string {
	return fmt.Sprintf("text=%s start=%d end=%d q_s=%v q_e=%v e=%v a=%v",
		string(t.Text),
		t.PosStart,
		t.PosEnd,
		t.IsQuoteStart,
		t.IsQuoteEnd,
		t.IsEllipsis,
		t.HasApostrophe,
	)
}

func (t *SentenceToken) Equals(compare *SentenceToken) bool {
	if string(t.Text) != string(compare.Text) {
		return false
	}
	if t.PosStart != compare.PosStart {
		return false
	}
	if t.PosEnd != compare.PosEnd {
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
