package book

type Interface interface {
	CreateBook(title string) (int64, error)
	UpdateBook(id int64, title string) error
	DeleteBook(id int64) error
	GetBooks(title string) ([]Book, error)
	GetBook(id int64) (*Book, error)

	CreateChapter(bookId int64, number int, blurURL string) (int64, error)
	UpdateChapter(id int64, number int, blurURL string) error
	DeleteChapters(bookId int64) error
	DeleteChapter(id int64) error
	GetChapters(bookId int64, number []int) ([]Chapter, error)
	GetChapter(chapterId int64) (*Chapter, error)

	CreatePage(page *Page) (int64, error)
	CreateManyPage(page []Page) ([]Page, error)
	UpdatePage(id int64, page *Page) error
	DeletePages(chapterId int64) error
	DeletePage(id int64) error
	GetPages(chapterIds []int64, withPanels bool) ([]Page, error)
	GetPage(pageId int64, withPanels bool) (*Page, error)

	CreatePanel(panel *Panel) (int64, error)
	CreateManyPanel(panel []Panel) ([]Panel, error)
	UpdatePanel(id int64, panel *Panel) error
	DeletePanels(pageId int64) error
	DeletePanel(id int64) error
	GetPanels(pageId int64) ([]Panel, error)
	GetPanel(panelId int64) (*Panel, error)
}
