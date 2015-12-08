package tokenize

type SentenceToken struct {
	Text          []rune
	PosStart      int
	PosEnd        int
	IsQuoteStart  bool
	IsQuoteEnd    bool
	IsEllipsis    bool
	HasApostrophe bool
}

func NewSentenceToken(str []rune, posStart, posEnd int) *SentenceToken {
	return &SentenceToken{
		Text:     str[posStart:posEnd],
		PosStart: posStart,
		PosEnd:   posEnd,
	}
}
