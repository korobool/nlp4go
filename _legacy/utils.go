package gonlp

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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

func ParseOntonotes(ontoPath string, outPath string) error {

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	walkFn := func(path string, f os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".parse") {
			file, _ := os.Open(path)
			defer file.Close()

			err := parseFile(file, writer)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = filepath.Walk(ontoPath, walkFn)
	if err != nil {
		return err
	}
	return nil
}

func parseFile(reader io.Reader, writer *bufio.Writer) error {

	sent := []byte{}
	re := regexp.MustCompile(`\s+`)

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if scanner.Text() != "" {
			sent = append(sent, scanner.Bytes()...)
			continue
		}
		if len(sent) > 0 {
			s := re.ReplaceAllLiteralString(string(sent), " ")
			fmt.Fprintln(writer, parseSent(s))
		}
		sent = []byte{}
	}
	err := scanner.Err()

	return err
}

func parseSent(tree string) string {

	byteStr := []byte("")
	re := regexp.MustCompile(`\(([^\s\(\)]+) ([^\s\(\)]+)\)`)

	match := re.FindAllStringSubmatch(tree, -1)
	for _, pair := range match {
		str := fmt.Sprintf("(%s %s)", pair[1], pair[2])
		byteStr = append(byteStr, []byte(str)...)
	}
	return string(byteStr)
}
