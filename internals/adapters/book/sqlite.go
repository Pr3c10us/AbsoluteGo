package book

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
)

type implementation struct {
	db *sql.DB
}

func NewBookImplementation(db *sql.DB) book.Interface {
	return &implementation{
		db: db,
	}
}
func (r *implementation) CreateBook(title string) error {
	_, err := sq.Insert("books").
		Columns("title").
		Values(title).
		RunWith(r.db).
		Exec()
	return err
}

func (r *implementation) UpdateBook(id int64, title string) error {
	res, err := sq.Update("books").
		Set("title", title).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (r *implementation) DeleteBook(id int64) error {
	res, err := sq.Delete("books").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (r *implementation) GetBooks(title string) ([]book.Book, error) {
	q := sq.Select("id", "title").From("books")
	if title != "" {
		q = q.Where(sq.Like{"title": fmt.Sprintf("%%%s%%", title)})
	}

	rows, err := q.RunWith(r.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []book.Book
	for rows.Next() {
		var b book.Book
		if err := rows.Scan(&b.Id, &b.Title); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

func (r *implementation) CreateChapter(bookId int64, number int, blurURL string) error {
	_, err := sq.Insert("chapters").
		Columns("book_id", "number", "blur_url").
		Values(bookId, number, blurURL).
		RunWith(r.db).
		Exec()
	return err
}

func (r *implementation) UpdateChapter(id int64, number int) error {
	res, err := sq.Update("chapters").
		Set("number", number).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (r *implementation) DeleteChapter(id int64) error {
	res, err := sq.Delete("chapters").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (r *implementation) GetChapters(bookId int64) ([]book.Chapter, error) {
	rows, err := sq.Select("id", "number", "book_id", "blur_url").
		From("chapters").
		Where(sq.Eq{"book_id": bookId}).
		OrderBy("number ASC").
		RunWith(r.db).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []book.Chapter
	for rows.Next() {
		var c book.Chapter
		if err = rows.Scan(&c.Id, &c.Number, &c.BookId, &c.BlurURL); err != nil {
			return nil, err
		}
		chapters = append(chapters, c)
	}
	return chapters, rows.Err()
}

func (r *implementation) CreatePage(page *book.Page) error {
	_, err := sq.Insert("pages").
		Columns("chapter_id", "url", "llm_url", "page_number").
		Values(page.ChapterId, page.URL, page.LLMURL, page.PageNumber).
		RunWith(r.db).
		Exec()
	return err
}

func (r *implementation) UpdatePage(id int64, page *book.Page) error {
	res, err := sq.Update("pages").
		Set("url", page.URL).
		Set("llm_url", page.LLMURL).
		Set("page_number", page.PageNumber).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (r *implementation) DeletePage(id int64) error {
	res, err := sq.Delete("pages").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (r *implementation) GetPages(chapterId int64) ([]book.Page, error) {
	rows, err := sq.Select("id", "chapter_id", "url", "llm_url", "page_number", "updated_at").
		From("pages").
		Where(sq.Eq{"chapter_id": chapterId}).
		OrderBy("page_number ASC").
		RunWith(r.db).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []book.Page
	for rows.Next() {
		var p book.Page
		if err := rows.Scan(&p.Id, &p.ChapterId, &p.URL, &p.LLMURL, &p.PageNumber, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, rows.Err()
}

func (r *implementation) CreatePanel(panel *book.Panel) error {
	_, err := sq.Insert("panels").
		Columns("page_id", "url", "llm_url", "panel_number").
		Values(panel.PageId, panel.URL, panel.LLMURL, panel.PanelNumber).
		RunWith(r.db).
		Exec()
	return err
}

func (r *implementation) UpdatePanel(id int64, panel *book.Panel) error {
	res, err := sq.Update("panels").
		Set("url", panel.URL).
		Set("llm_url", panel.LLMURL).
		Set("panel_number", panel.PanelNumber).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (r *implementation) DeletePanel(id int64) error {
	res, err := sq.Delete("panels").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (r *implementation) GetPanels(pageId int64) ([]book.Panel, error) {
	rows, err := sq.Select("id", "page_id", "url", "llm_url", "panel_number", "updated_at").
		From("panels").
		Where(sq.Eq{"page_id": pageId}).
		OrderBy("panel_number ASC").
		RunWith(r.db).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var panels []book.Panel
	for rows.Next() {
		var p book.Panel
		if err := rows.Scan(&p.Id, &p.PageId, &p.URL, &p.LLMURL, &p.PanelNumber, &p.UpdatedAt); err != nil {
			return nil, err
		}
		panels = append(panels, p)
	}
	return panels, rows.Err()
}

func assertRowAffected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
