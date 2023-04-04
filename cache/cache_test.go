package cache_test

import (
	"test_H5B/cache"
	"test_H5B/output"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testerBuild struct {
	Input  inputBuild
	Output outputBuild
}

type inputBuild struct {
	Bytes []byte
}

type outputBuild struct {
	Out map[string]*output.Search
	Err error
}

var testBuild = map[string]testerBuild{
	"ok empty": {
		Input: inputBuild{
			Bytes: []byte("..."),
		},
		Output: outputBuild{
			Out: map[string]*output.Search{},
			Err: nil,
		},
	},
	"ok with single line one occurence": {
		Input: inputBuild{
			Bytes: []byte("Lorem"),
		},
		Output: outputBuild{
			Out: map[string]*output.Search{
				"Lorem": {
					WordFound:       true,
					NumOccurrences:  1,
					LineOccurrences: []int{1},
				},
			},
			Err: nil,
		},
	},
	"ok with single line one occurence with unknown post key": {
		Input: inputBuild{
			Bytes: []byte("Lorem."),
		},
		Output: outputBuild{
			Out: map[string]*output.Search{
				"Lorem": {
					WordFound:       true,
					NumOccurrences:  1,
					LineOccurrences: []int{1},
				},
			},
			Err: nil,
		},
	},
	"ok with single line one occurence with unknown prior key": {
		Input: inputBuild{
			Bytes: []byte("!Lorem"),
		},
		Output: outputBuild{
			Out: map[string]*output.Search{
				"Lorem": {
					WordFound:       true,
					NumOccurrences:  1,
					LineOccurrences: []int{1},
				},
			},
			Err: nil,
		},
	},
	"ok with multine line one occurence with unknown keys": {
		Input: inputBuild{
			Bytes: []byte("!  Lorem   .   \n"),
		},
		Output: outputBuild{
			Out: map[string]*output.Search{
				"Lorem": {
					WordFound:       true,
					NumOccurrences:  1,
					LineOccurrences: []int{1},
				},
			},
			Err: nil,
		},
	},
	"ok with multiline line one occurence with unknown keys": {
		Input: inputBuild{
			Bytes: []byte("\n! \n  Lorem   . \n   \n"),
		},
		Output: outputBuild{
			Out: map[string]*output.Search{
				"Lorem": {
					WordFound:       true,
					NumOccurrences:  1,
					LineOccurrences: []int{3},
				},
			},
			Err: nil,
		},
	},
	"ok with multiline line multiple occurence with unknown keys": {
		Input: inputBuild{
			Bytes: []byte("\n!Lorem \n  Lorem   . \n   Lorem"),
		},
		Output: outputBuild{
			Out: map[string]*output.Search{
				"Lorem": {
					WordFound:       true,
					NumOccurrences:  3,
					LineOccurrences: []int{2, 3, 4},
				},
			},
			Err: nil,
		},
	},
	"ok with multiline line multiple occurence with unknown keys with same line": {
		Input: inputBuild{
			Bytes: []byte("\n!Lorem.Lorem ? \n  Lorem Lorem Lorem  . \n   Lorem"),
		},
		Output: outputBuild{
			Out: map[string]*output.Search{
				"Lorem": {
					WordFound:       true,
					NumOccurrences:  6,
					LineOccurrences: []int{2, 3, 4},
				},
			},
			Err: nil,
		},
	},
	"ok with multiple keys with multiline line multiple occurence with unknown keys with same line": {
		Input: inputBuild{
			Bytes: []byte("\n!Lorem.Lorem lorem ipsum? \n  Lorem Lorem Lorem  . \n   lorem: Lorem"),
		},
		Output: outputBuild{
			Out: map[string]*output.Search{
				"Lorem": {
					WordFound:       true,
					NumOccurrences:  6,
					LineOccurrences: []int{2, 3, 4},
				},
				"lorem": {
					WordFound:       true,
					NumOccurrences:  2,
					LineOccurrences: []int{2, 4},
				},
				"ipsum": {
					WordFound:       true,
					NumOccurrences:  1,
					LineOccurrences: []int{2},
				},
			},
			Err: nil,
		},
	},
}

func TestBuild(t *testing.T) {
	for k, v := range testBuild {
		t.Run(k, func(t *testing.T) {
			out, err := cache.Build(v.Input.Bytes)
			if v.Output.Err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, v.Output.Out, out)
			}

		})
	}
}
