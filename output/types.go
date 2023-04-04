package output

type Search struct {
	WordFound       bool  `json:"wordFound"`
	NumOccurrences  int   `json:"numOccurrences"`
	LineOccurrences []int `json:"lineOccurrences"`
}

type Error struct {
	Error string `json:"error"`
}
