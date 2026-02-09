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
func (i *implementation) CreateBook(title string) (int64, error) {
	res, err := sq.Insert("books").
		Columns("title").
		Values(title).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) UpdateBook(id int64, title string) error {
	q := sq.Update("books").Where(sq.Eq{"id": id})
	if title != "" {
		q = q.Set("title", title)
	}
	res, err := q.RunWith(i.db).Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) DeleteBook(id int64) error {
	res, err := sq.Delete("books").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetBooks(title string) ([]book.Book, error) {
	q := sq.Select("id", "title").From("books")
	if title != "" {
		q = q.Where(sq.Like{"title": fmt.Sprintf("%%%s%%", title)})
	}

	rows, err := q.RunWith(i.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []book.Book{}
	for rows.Next() {
		var b book.Book
		if err := rows.Scan(&b.Id, &b.Title); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

func (i *implementation) GetBook(id int64) (*book.Book, error) {
	var b book.Book
	err := sq.Select("id", "title").
		From("books").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		QueryRow().
		Scan(&b.Id, &b.Title)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (i *implementation) CreateChapter(bookId int64, number int, blurURL string) (int64, error) {
	res, err := sq.Insert("chapters").
		Columns("book_id", "number", "blur_url").
		Values(bookId, number, blurURL).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) UpdateChapter(id int64, number int, blurURL string) error {
	q := sq.Update("chapters").Where(sq.Eq{"id": id})
	if number > 0 {
		q = q.Set("number", number)
	}
	if blurURL != "" {
		q = q.Set("blur_url", blurURL)
	}
	res, err := q.RunWith(i.db).Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) DeleteChapters(bookId int64) error {
	_, err := sq.Delete("chapters").
		Where(sq.Eq{"book_id": bookId}).
		RunWith(i.db).
		Exec()
	return err
}

func (i *implementation) DeleteChapter(id int64) error {
	res, err := sq.Delete("chapters").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetChapters(bookId int64, number int) ([]book.Chapter, error) {
	q := sq.Select("id", "number", "book_id", "blur_url").
		From("chapters").
		Where(sq.Eq{"book_id": bookId}).
		OrderBy("number ASC")
	if number > 0 {
		q = q.Where(sq.Eq{"number": number})
	}

	rows, err := q.RunWith(i.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chapters := []book.Chapter{}
	for rows.Next() {
		var c book.Chapter
		if err = rows.Scan(&c.Id, &c.Number, &c.BookId, &c.BlurURL); err != nil {
			return nil, err
		}
		chapters = append(chapters, c)
	}
	return chapters, rows.Err()
}

func (i *implementation) GetChapter(chapterId int64) (*book.Chapter, error) {
	var c book.Chapter
	err := sq.Select("id", "number", "book_id", "blur_url").
		From("chapters").
		Where(sq.Eq{"id": chapterId}).
		RunWith(i.db).
		QueryRow().
		Scan(&c.Id, &c.Number, &c.BookId, &c.BlurURL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (i *implementation) CreateManyPage(pages []book.Page) ([]book.Page, error) {
	tx, err := i.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result := make([]book.Page, len(pages))
	for idx, p := range pages {
		res, err := sq.Insert("pages").
			Columns("chapter_id", "url", "llm_url", "mime", "page_number").
			Values(p.ChapterId, p.URL, p.LLMURL, p.MIME, p.PageNumber).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		p.Id = id
		result[idx] = p
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (i *implementation) CreatePage(page *book.Page) (int64, error) {
	res, err := sq.Insert("pages").
		Columns("chapter_id", "url", "llm_url", "mime", "page_number").
		Values(page.ChapterId, page.URL, page.LLMURL, page.MIME, page.PageNumber).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) UpdatePage(id int64, page *book.Page) error {
	q := sq.Update("pages").Where(sq.Eq{"id": id})
	if page.URL != nil {
		q = q.Set("url", page.URL)
	}
	if page.LLMURL != nil {
		q = q.Set("llm_url", page.LLMURL)
	}
	if page.MIME != nil {
		q = q.Set("mime", page.MIME)
	}
	if page.PageNumber > 0 {
		q = q.Set("page_number", page.PageNumber)
	}
	res, err := q.RunWith(i.db).Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) DeletePages(chapterId int64) error {
	_, err := sq.Delete("pages").
		Where(sq.Eq{"chapter_id": chapterId}).
		RunWith(i.db).
		Exec()
	return err
}

func (i *implementation) DeletePage(id int64) error {
	res, err := sq.Delete("pages").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetPages(chapterId int64) ([]book.Page, error) {
	rows, err := sq.Select("id", "chapter_id", "url", "llm_url", "mime", "page_number", "updated_at").
		From("pages").
		Where(sq.Eq{"chapter_id": chapterId}).
		OrderBy("page_number ASC").
		RunWith(i.db).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pages := []book.Page{}
	for rows.Next() {
		var p book.Page
		if err := rows.Scan(&p.Id, &p.ChapterId, &p.URL, &p.LLMURL, &p.MIME, &p.PageNumber, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, rows.Err()
}

func (i *implementation) GetPage(pageId int64) (*book.Page, error) {
	var p book.Page
	err := sq.Select("id", "chapter_id", "url", "llm_url", "mime", "page_number", "updated_at").
		From("pages").
		Where(sq.Eq{"id": pageId}).
		RunWith(i.db).
		QueryRow().
		Scan(&p.Id, &p.ChapterId, &p.URL, &p.LLMURL, &p.MIME, &p.PageNumber, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (i *implementation) CreateManyPanel(panels []book.Panel) ([]book.Panel, error) {
	tx, err := i.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result := make([]book.Panel, len(panels))
	for idx, p := range panels {
		res, err := sq.Insert("panels").
			Columns("page_id", "url", "panel_number").
			Values(p.PageId, p.URL, p.PanelNumber).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		p.Id = id
		result[idx] = p
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (i *implementation) CreatePanel(panel *book.Panel) (int64, error) {
	res, err := sq.Insert("panels").
		Columns("page_id", "url", "panel_number").
		Values(panel.PageId, panel.URL, panel.PanelNumber).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) UpdatePanel(id int64, panel *book.Panel) error {
	q := sq.Update("panels").Where(sq.Eq{"id": id})
	if panel.URL != nil {
		q = q.Set("url", panel.URL)
	}
	if panel.PanelNumber > 0 {
		q = q.Set("panel_number", panel.PanelNumber)
	}
	res, err := q.RunWith(i.db).Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) DeletePanels(pageId int64) error {
	_, err := sq.Delete("panels").
		Where(sq.Eq{"page_id": pageId}).
		RunWith(i.db).
		Exec()
	return err
}

func (i *implementation) DeletePanel(id int64) error {
	res, err := sq.Delete("panels").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetPanels(pageId int64) ([]book.Panel, error) {
	rows, err := sq.Select("id", "page_id", "url", "panel_number", "updated_at").
		From("panels").
		Where(sq.Eq{"page_id": pageId}).
		OrderBy("panel_number ASC").
		RunWith(i.db).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	panels := []book.Panel{}
	for rows.Next() {
		var p book.Panel
		if err := rows.Scan(&p.Id, &p.PageId, &p.URL, &p.PanelNumber, &p.UpdatedAt); err != nil {
			return nil, err
		}
		panels = append(panels, p)
	}
	return panels, rows.Err()
}

func (i *implementation) GetPanel(panelId int64) (*book.Panel, error) {
	var p book.Panel
	err := sq.Select("id", "page_id", "url", "panel_number", "updated_at").
		From("panels").
		Where(sq.Eq{"id": panelId}).
		RunWith(i.db).
		QueryRow().
		Scan(&p.Id, &p.PageId, &p.URL, &p.PanelNumber, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
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
