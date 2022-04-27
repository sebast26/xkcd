package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Comic represents single XKCD comic strip
type Comic struct {
	URL        string `json:"img"`
	Transcript string `json:"transcript"`
}

func matchesTerm(comic *Comic, term string) bool {
	return strings.Contains(comic.Transcript, term)
}

func parse(handle io.Reader) (*Comic, error) {
	var comic Comic
	err := json.NewDecoder(handle).Decode(&comic)
	if err != nil {
		return nil, fmt.Errorf("error decoding: %v", err)
	}
	return &comic, nil
}
