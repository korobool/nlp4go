package nlp_string

import (
	"bytes"
	"regexp"
	"unicode/utf8"
)

type noRuneError struct {
}

func (e noRuneError) Error() string {
	return "no next rune in string"
}

type String struct {
	runes  []rune
	length int
	ridx   int
}

func NewString(orig string) *String {
	rc := utf8.RuneCountInString(orig)
	runes := make([]rune, rc)
	rune_pos := 0
	pos := 0
	for {
		r, size := utf8.DecodeRuneInString(orig[pos:])
		if size == 0 {
			break
		}
		pos += size
		runes[rune_pos] = r
		rune_pos += 1
	}
	s := String{runes: runes, length: int(rc)}
	return &s
}

func (s *String) Length() int {
	return s.length
}

func (s *String) String() string {
	var size int
	var ret_bytes bytes.Buffer
	buf := make([]byte, 4)

	for _, r := range s.runes {
		size = utf8.EncodeRune(buf, r)
		ret_bytes.Write(buf[:size])
	}
	return ret_bytes.String()
}

func (s *String) Substring(from, to int) *String {
	if from < 0 || from > to || to > s.length {
		return nil
	}
	runes := s.runes[from:to]
	new_s := String{runes: runes, length: to - from}
	return &new_s
}

func (s *String) ReadRune() (r rune, size int, err error) {
	if s.ridx == s.length {
		return 0, 0, noRuneError{}
	}
	r = s.runes[s.ridx]
	err = nil
	s.ridx += 1
	// padding to 4 bytes per rune
	// expecting offset returned from regexp.FindReaderIndex equal to rune_offset*4
	size = 4
	return
}

func (s *String) normalizeRegexpLoc(runes_offset, from, to int) (loc []int) {
	from_rune := runes_offset + (from / 4)
	to_rune := runes_offset + (to / 4)
	loc = []int{from_rune, to_rune}
	return loc
}

func (s *String) FindAll(re *regexp.Regexp) [][]int {
	s.ridx = 0
	var ret_locs [][]int
	var cur_ridx int
	for {
		cur_ridx = s.ridx
		b_loc := re.FindReaderIndex(s)
		if b_loc == nil {
			s.ridx = 0
			break
		}
		ret_locs = append(ret_locs, s.normalizeRegexpLoc(cur_ridx, b_loc[0], b_loc[1]))

	}
	return ret_locs
}

func (s *String) FindFirst(re *regexp.Regexp) (loc []int) {
	s.ridx = 0
	b_loc := re.FindReaderIndex(s)
	if b_loc == nil {
		s.ridx = 0
		return nil
	}
	s.ridx = 0
	return s.normalizeRegexpLoc(0, b_loc[0], b_loc[1])
}

func (s *String) Match(re *regexp.Regexp) bool {
	s.ridx = 0
	matched := re.MatchReader(s)
	s.ridx = 0
	return matched
}

//func (s *String) Replace(re *regexp.Regexp, newText string) {
//
//}
