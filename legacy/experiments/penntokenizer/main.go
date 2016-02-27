package main

import (
	"fmt"
	"regexp"
)

func main() {

	res := []*regexp.Regexp{}

	text := "They'll save 5% and 10$ invest more."

	res = append(res, regexp.MustCompile("(``)"))
	res = append(res, regexp.MustCompile(`([:,])[^\d]`))
	res = append(res, regexp.MustCompile(`(\.\.\.)`))
	res = append(res, regexp.MustCompile(`([;@#$%&?!])`))
	res = append(res, regexp.MustCompile(`[^\.](\.)[\]\)}>"']*\s*$`))
	res = append(res, regexp.MustCompile(`[^'](') `))
	res = append(res, regexp.MustCompile(`([\]\[\(\)\{\}\<\>])`))
	res = append(res, regexp.MustCompile(`(--)`))
	res = append(res, regexp.MustCompile(`(")`))
	res = append(res, regexp.MustCompile(`([^' ]+)('[sS]|'[mM]|'[dD]|'ll|'LL|'re|'RE|'ve|'VE|n't|N'T|') `))
	res = append(res, regexp.MustCompile(`(?i)\b(can)(not)\b`))
	res = append(res, regexp.MustCompile(`(?i)\b(d)('ye)\b`))
	res = append(res, regexp.MustCompile(`(?i)\b(gim)(me)\b`))
	res = append(res, regexp.MustCompile(`(?i)\b(gon)(na)\b`))
	res = append(res, regexp.MustCompile(`(?i)\b(got)(ta)\b`))
	res = append(res, regexp.MustCompile(`(?i)\b(lem)(me)\b`))
	res = append(res, regexp.MustCompile(`(?i)\b(mor)('n)\b`))
	res = append(res, regexp.MustCompile(`(?i)\b(wan)(na) `))
	res = append(res, regexp.MustCompile(`(?i) ('t)(is)\b`))
	res = append(res, regexp.MustCompile(`(?i) ('t)(was)\b`))

	for i, re := range res {
		matches := re.FindAllStringSubmatchIndex(text, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				for indx := 2; indx < len(match); indx += 2 {
					cors := match[indx : indx+2]
					fmt.Println(i, cors, text[cors[0]:cors[1]])
				}
			}
		}
	}
}
