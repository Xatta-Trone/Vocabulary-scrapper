package model

import "strconv"

type Word struct {
	ID         int    `json:"id,omitempty"`
	Word       string `json:"word"`
	Definition string `json:"definition,omitempty"`
}

type Options struct {
	URL           string
	OUTPUT        string
	NO_DEFINITION bool
	NO_ID         bool
	ONLY_WORD     bool
	ONLY_CSV      bool
	ONLY_JSON     bool
}

func (w Word) BuildCSV() []string {

	if w.ID == 0 {
		return []string{w.Word, w.Definition}
	}
	return []string{strconv.Itoa(w.ID), w.Word, w.Definition}

}
