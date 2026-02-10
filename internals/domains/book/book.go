package book

import "time"

type Book struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

type Chapter struct {
	Id      int64  `json:"id"`
	Number  int    `json:"number"`
	BookId  int64  `json:"bookId"`
	BlurURL string `json:"blurURL"`
}

type Page struct {
	Id         int64     `json:"id"`
	ChapterId  int64     `json:"chapterId"`
	URL        *string   `json:"url"`
	LLMURL     *string   `json:"llmurl"`
	MIME       *string   `json:"mime"`
	PageNumber int       `json:"pageNumber"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Panels     []Panel   `json:"panels"`
}

type Panel struct {
	Id          int64     `json:"id"`
	PageId      int64     `json:"pageId"`
	URL         *string   `json:"url"`
	PanelNumber int       `json:"panelNumber"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
