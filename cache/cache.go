package cache

import (
	"bytes"
	"io"
	"strings"
	"test_H5B/output"
	"unicode"
)

// Builds from byte array to store all possible results in memory
func Build(content []byte) (map[string]*output.Search, error) {
	out := make(map[string]*output.Search)
	reader := bytes.NewReader(content)
	currentLine := 1
	word := &strings.Builder{}
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			word.WriteRune(r)
			continue
		}
		saveWord(word, currentLine, out)
		if r == '\n' {
			currentLine += 1
		}
	}
	return out, nil
}

func saveWord(word *strings.Builder, currentLine int, out map[string]*output.Search) {
	currentWord := word.String()
	if currentWord == "" {
		return
	}
	if _, found := out[currentWord]; found {
		out[currentWord].NumOccurrences += 1
		// Info: condition is safe since LineOccurences will be always at a minimum of one due to else condition
		if out[currentWord].LineOccurrences[len(out[currentWord].LineOccurrences)-1] != currentLine {
			out[currentWord].LineOccurrences = append(out[currentWord].LineOccurrences, currentLine)
		}
	} else {
		out[currentWord] = &output.Search{
			WordFound:       true,
			NumOccurrences:  1,
			LineOccurrences: []int{currentLine},
		}
	}
	word.Reset()
}
