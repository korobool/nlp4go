package strings

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZйцукенгшщзхъфывапролджэячсмитьбю")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func BenchmarkNewString(b *testing.B) {
	s := RandStringRunes(1000000)
	for i := 0; i < b.N; i++ {
		_ = NewString(s)
	}
}

func BenchmarkSubstring(b *testing.B) {
	s := RandStringRunes(1000000)
	str := NewString(s)
	for i := 0; i < b.N; i++ {
		_ = str.Substring(i%1000000, 1000000)
	}
}

func BenchmarkToString(b *testing.B) {
	s := RandStringRunes(1000000)
	str := NewString(s)
	for i := 0; i < b.N; i++ {
		_ = str.String()
	}
}

func BenchmarkFindAll_Many(b *testing.B) {
	s := RandStringRunes(1000000)
	str := NewString(s)
	re, _ := regexp.Compile("a")
	for i := 0; i < b.N; i++ {
		_ = str.FindAll(re)
	}
}

func BenchmarkFindAll_NoOne(b *testing.B) {
	s := RandStringRunes(1000000)
	str := NewString(s)
	re, _ := regexp.Compile("HOHOHO not found this")
	for i := 0; i < b.N; i++ {
		_ = str.FindAll(re)
	}
}

func BenchmarkReplace(b *testing.B) {
	s := RandStringRunes(1000000)
	str := NewString(s)
	re, _ := regexp.Compile("A")
	for i := 0; i < b.N; i++ {
		_, _ = str.Replace(re, "$")
		//fmt.Println("Replaced: ", cnt)
	}
}

func BenchmarkReplaceNotFound(b *testing.B) {
	s := RandStringRunes(1000000)
	str := NewString(s)
	re, _ := regexp.Compile("NOT_FOUND")
	for i := 0; i < b.N; i++ {
		_, _ = str.Replace(re, "hohoho")
	}
}

func TestBasic(t *testing.T) {
	init_str := "String with unicode: こんに"
	str := NewString(init_str)
	assert.Equal(t, len(init_str), 30)
	assert.Equal(t, str.Length(), 24)
	assert.Equal(t, str.String(), init_str)

	sub := str.Substring(21, 23)
	assert.Equal(t, sub.Length(), 2)
	assert.Equal(t, sub.String(), "こん")
	sub = str.Substring(4, 3)
	assert.Nil(t, sub)
	sub = str.Substring(-4, 3)
	assert.Nil(t, sub)
	sub = str.Substring(1, 30)
	assert.Nil(t, sub)
	sub = str.Substring(28, 12)
	assert.Nil(t, sub)

	re1, _ := regexp.Compile("with.*")
	re2, _ := regexp.Compile("i")
	re3, _ := regexp.Compile("ん.*")
	re4, _ := regexp.Compile("string")
	matched := str.Match(re1)
	assert.Equal(t, matched, true)
	matched = str.Match(re2)
	assert.Equal(t, matched, true)
	matched = str.Match(re3)
	assert.Equal(t, matched, true)
	matched = str.Match(re4)
	assert.Equal(t, matched, false)

	loc := str.FindFirst(re1)
	assert.Equal(t, loc, []int{7, 24})
	sub = str.Substring(loc[0], loc[1])
	assert.Equal(t, sub.String(), "with unicode: こんに")
	loc = str.FindFirst(re4)
	assert.Nil(t, nil)

	locs := str.FindAll(re2)
	assert.Equal(t, len(locs), 3)
	assert.Equal(t, locs[0], []int{3, 4})
	assert.Equal(t, locs[1], []int{8, 9})
	assert.Equal(t, locs[2], []int{14, 15})

	locs = str.FindAll(re3)
	assert.Equal(t, len(locs), 1)
	assert.Equal(t, locs[0], []int{22, 24})

	locs = str.FindAll(re4)
	assert.Equal(t, len(locs), 0)

	newS, cnt := str.Replace(re2, "Ё$")
	assert.Equal(t, cnt, 3)
	assert.Equal(t, newS.Length(), 27)
	assert.Equal(t, newS.String(), "StrЁ$ng wЁ$th unЁ$code: こんに")

	newS, cnt = str.Replace(re3, "")
	assert.Equal(t, cnt, 1)
	assert.Equal(t, newS.String(), "String with unicode: こ")

	newS, cnt = str.Replace(re4, "HELLO")
	assert.Equal(t, cnt, 0)
	assert.Equal(t, newS.String(), str.String())
}
