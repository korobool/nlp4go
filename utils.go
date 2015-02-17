package gonlp

import (
	"math"
	"sort"
)

type Pair struct {
	Key   string
	Value float64
}

type PairList []Pair

func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int      { return len(p) }
func (p PairList) Less(i, j int) bool {
	if p[i].Value == p[j].Value {
		return p[i].Key > p[j].Key
	}
	return p[i].Value > p[j].Value
}

func SortMapKV(m map[string]float64) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i += 1
	}
	sort.Sort(p)
	return p
}

func MaxScore(m map[string]float64) string {
	var maxKey string
	var maxVal float64

	for k, v := range m {
		switch {
		case v > maxVal:
			maxKey = k
			maxVal = v
		case v == maxVal:
			if k < maxKey {
				maxKey = k
				maxVal = v
			}
		}
	}
	return maxKey

}

// https://gist.github.com/pelegm/c48cff315cd223f7cf7b
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(roundOn, val)
	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
