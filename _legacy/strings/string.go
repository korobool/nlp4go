package strings

import (
	"errors"
	"unicode/utf8"
)

var ErrIndex = errors.New("index error")

type String struct {
	String       string
	Bytes        []byte
	Runes        []rune
	bytesToRunes []int
	runesToBytes []int
}

func (s String) Clone() *String {

	c := String{
		String:       s.String,
		Bytes:        make([]byte, len(s.Bytes)),
		Runes:        make([]rune, len(s.Runes)),
		bytesToRunes: make([]int, len(s.bytesToRunes)),
		runesToBytes: make([]int, len(s.runesToBytes)),
	}
	copy(c.Bytes, s.Bytes)
	copy(c.Runes, s.Runes)
	copy(c.bytesToRunes, s.bytesToRunes)
	copy(c.runesToBytes, s.runesToBytes)

	return &c
}

func (s *String) ByteToRuneRange(from, to int) (int, int, error) {
	if from > to || from < 0 || to < 0 {
		return 0, 0, ErrIndex
	}
	if to > len(s.bytesToRunes) {
		return 0, 0, ErrIndex
	}

	runeFrom := s.bytesToRunes[from]
	runeTo := s.bytesToRunes[to]

	if (from != 0 && runeFrom == 0) || (to != 0 && runeTo == 0) {
		return 0, 0, ErrIndex
	}
	return runeFrom, runeTo, nil
}

func FromString(str string) *String {
	byteStr := []byte(str)
	s := prepareFromBytes(byteStr)
	s.String = str
	return s
}

func FromBytes(byteStr []byte) *String {
	s := prepareFromBytes(byteStr)
	s.String = string(byteStr)
	return s
}

func prepareFromBytes(byteStr []byte) *String {

	runesSlice := make([]rune, len(byteStr))
	bytesToRunes := make([]int, len(byteStr))
	runesToBytes := make([]int, len(byteStr))

	bpos, rpos := 0, 0
	for w := 0; bpos < len(byteStr); bpos += w {
		r, width := utf8.DecodeRune(byteStr[bpos:])

		runesSlice[rpos] = r
		runesToBytes[rpos] = bpos
		bytesToRunes[bpos] = rpos

		w = width
		rpos++
	}
	s := &String{
		Bytes:        byteStr,
		bytesToRunes: bytesToRunes,
	}

	if rpos < cap(runesSlice) {
		s.Runes = make([]rune, len(runesSlice[:rpos]))
		s.runesToBytes = make([]int, len(runesToBytes[:rpos]))

		copy(s.Runes, runesSlice[:rpos])
		copy(s.runesToBytes, runesToBytes[:rpos])

	} else {
		s.Runes = runesSlice
		s.runesToBytes = runesToBytes
	}
	return s
}
