package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
)

const (
	fileName = "Latin-Lipsum.txt"
	port     = 62626
)

type output struct {
	WordFound       bool  `json:"wordFound"`
	NumOccurrences  int   `json:"numOccurrences"`
	LineOccurrences []int `json:"lineOccurrences"`
}

type errorOutput struct {
	Error string `json:"error"`
}

var memCache map[string]*output

func init() {
	b, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	if memCache, err = buildCache(b); err != nil {
		panic(err)
	}
}

func search(c echo.Context) error {
	var searchWord string
	out := &output{
		LineOccurrences: []int{},
	}
	if err := echo.PathParamsBinder(c).String("searchWord", &searchWord).BindError(); err != nil {
		// Info: On fails should be 500 => this error should never occur
		return c.JSON(http.StatusInternalServerError, errorOutput{Error: err.Error()})
	}
	if _, found := memCache[searchWord]; found {
		out = memCache[searchWord]
	}
	return c.JSON(http.StatusOK, out)
}

func main() {
	e := echo.New()
	searchPrefix := e.Group("/api/v0.1/search")
	searchPrefix.GET("/:searchWord", search)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

// Builds from byte array to store all possible results in memory
func buildCache(content []byte) (map[string]*output, error) {
	out := make(map[string]*output)
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

func saveWord(word *strings.Builder, currentLine int, out map[string]*output) {
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
		out[currentWord] = &output{
			WordFound:       true,
			NumOccurrences:  1,
			LineOccurrences: []int{currentLine},
		}
	}
	word.Reset()
}
