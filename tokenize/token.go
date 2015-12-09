package tokenize

type SentenceToken struct {
	Text          []rune `json:"text"`
	PosStart      int    `json:"pos_start"`
	PosEnd        int    `json:"pos_end"`
	IsQuoteStart  bool   `json:"is_quote_start"`
	IsQuoteEnd    bool   `json:"is_quote_end"`
	IsEllipsis    bool   `json:"is_ellipsis"`
	HasApostrophe bool   `json:"has_apostrophe"`
}

func (t *SentenceToken) String() string {
	return string(t.Text)
}

func NewSentenceToken(str []rune, posStart, posEnd int) *SentenceToken {
	return &SentenceToken{
		Text:     str[posStart:posEnd],
		PosStart: posStart,
		PosEnd:   posEnd,
	}
}
