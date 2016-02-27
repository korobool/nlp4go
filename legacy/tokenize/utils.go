package tokenize

import "unicode/utf8"

func byteToRunePosition(s string, left, right int) (int, int) {
	boundL, boundR := 0, 0
	bpos, rpos := 0, 0

	for b := 0; bpos < len(s); bpos += b {
		_, size := utf8.DecodeRuneInString(s[bpos:])

		if bpos == left {
			boundL = rpos
		} else if bpos == right-size {
			boundR = rpos + size
		}
		b = size
		rpos++
	}
	return boundL, boundR
}
