package book

type Interface interface {
	CreateBook(title string) error
	UpdateBook(id int64, title string) error
	DeleteBook(id int64) error
	GetBooks(title string) ([]Book, error)

	CreateChapter(bookId int64, number int, blurURL string) error
	UpdateChapter(id int64, number int) error
	DeleteChapter(id int64) error
	GetChapters(bookId int64) ([]Chapter, error)

	CreatePage(page *Page) error
	UpdatePage(id int64, page *Page) error
	DeletePage(id int64) error
	GetPages(chapterId int64) ([]Page, error)

	CreatePanel(panel *Panel) error
	UpdatePanel(id int64, panel *Panel) error
	DeletePanel(id int64) error
	GetPanels(pageId int64) ([]Panel, error)
}
