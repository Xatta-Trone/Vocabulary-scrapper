package model

import "strconv"

type Word struct {
	ID         int    `json:"id,omitempty"`
	Word       string `json:"word"`
	Definition string `json:"definition,omitempty"`
	Group      int    `json:"group,omitempty"`
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

type ResponseModel struct {
	FolderURL string                `json:"folder_url"`
	Sets      []SingleResponseModel `json:"sets"`
}

type SingleResponseModel struct {
	Title   string   `json:"title"`
	GroupId int      `json:"group_id"`
	URL     string   `json:"url"`
	Words   []string `json:"words"`
}

type QuizletFolder struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
}

func (w Word) BuildCSV() []string {

	if w.ID == 0 {
		return []string{w.Word, w.Definition}
	}
	return []string{strconv.Itoa(w.ID), strconv.Itoa(w.Group), w.Word, w.Definition}

}
