package script

import (
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/packages/utils"
)

type Script struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Content  *string `json:"content"`
	BookId   int64   `json:"bookId"`
	Chapters []int   `json:"chapters"`
}

type Query struct {
	BookId  int64
	Name    string
	Ids     []int64
	Chapter int
}

type Split struct {
	Id              int64         `json:"id"`
	ScriptId        int64         `json:"scriptId"`
	Content         *string       `json:"content"`
	PreviousContent *string       `json:"previousContent"`
	PanelId         *int64        `json:"panelId"`
	Effect          *utils.Effect `json:"effect"`
	AudioURL        *string       `json:"audioURL"`
	AudioDuration   *float64      `json:"audioDuration"`
	VideoURL        *string       `json:"videoURL"`
	Panel           book.Panel    `json:"panel"`
}
