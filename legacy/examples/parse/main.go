package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	walkFn := func(path string, f os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".parse") {
			file, _ := os.Open(path)
			defer file.Close()

			err := printSent(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		return nil
	}

	ontoPath := flag.String("p", "ontonotes", "ontonotes dir path")
	flag.Parse()

	filepath.Walk(*ontoPath, walkFn)

	// filepath.Walk("/home/demyan/ontonotes", walkFn)

}

func printSent(reader io.Reader) error {

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
			fmt.Println(parseSent(s))
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
