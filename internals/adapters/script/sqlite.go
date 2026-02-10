package script

import (
	"database/sql"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
)

type implementation struct {
	db *sql.DB
}

func NewScriptImplementation(db *sql.DB) script.Interface {
	return &implementation{
		db: db,
	}
}

func (i *implementation) CreateScript(s *script.Script) (int64, error) {
	chaptersJSON, err := json.Marshal(s.Chapters)
	if err != nil {
		return 0, err
	}
	res, err := sq.Insert("scripts").
		Columns("name", "content", "book_id", "chapters").
		Values(s.Name, s.Content, s.BookId, string(chaptersJSON)).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) UpdateScript(id int64, s *script.Script) error {
	q := sq.Update("scripts").Where(sq.Eq{"id": id})
	if s.Name != "" {
		q = q.Set("name", s.Name)
	}
	if s.Content != nil {
		q = q.Set("content", s.Content)
	}
	if s.Chapters != nil {
		chaptersJSON, err := json.Marshal(s.Chapters)
		if err != nil {
			return err
		}
		q = q.Set("chapters", string(chaptersJSON))
	}
	res, err := q.RunWith(i.db).Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) DeleteScript(id int64) error {
	res, err := sq.Delete("scripts").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetScripts(bookId int64, name string) ([]script.Script, error) {
	q := sq.Select("id", "name", "content", "book_id", "chapters").
		From("scripts")
	if bookId > 0 {
		q = q.Where(sq.Eq{"book_id": bookId})
	}
	if name != "" {
		q = q.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", name)})
	}

	rows, err := q.RunWith(i.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scripts := []script.Script{}
	for rows.Next() {
		var s script.Script
		var chaptersJSON string
		if err := rows.Scan(&s.Id, &s.Name, &s.Content, &s.BookId, &chaptersJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(chaptersJSON), &s.Chapters); err != nil {
			return nil, err
		}
		scripts = append(scripts, s)
	}
	return scripts, rows.Err()
}

func (i *implementation) GetScript(id int64) (*script.Script, error) {
	var s script.Script
	var chaptersJSON string
	err := sq.Select("id", "name", "content", "book_id", "chapters").
		From("scripts").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		QueryRow().
		Scan(&s.Id, &s.Name, &s.Content, &s.BookId, &chaptersJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(chaptersJSON), &s.Chapters); err != nil {
		return nil, err
	}
	return &s, nil
}

func (i *implementation) CreateSplit(split *script.Split) (int64, error) {
	res, err := sq.Insert("splits").
		Columns("script_id", "content", "panel_id", "effect").
		Values(split.ScriptId, split.Content, split.PanelId, split.Effect).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) CreateManySplit(splits []script.Split) ([]script.Split, error) {
	tx, err := i.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result := make([]script.Split, len(splits))
	for idx, s := range splits {
		res, err := sq.Insert("splits").
			Columns("script_id", "content", "panel_id", "effect").
			Values(s.ScriptId, s.Content, s.PanelId, s.Effect).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		s.Id = id
		result[idx] = s
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (i *implementation) UpdateSplit(id int64, split *script.Split) error {
	q := sq.Update("splits").Where(sq.Eq{"id": id})
	if split.Content != nil {
		q = q.Set("content", split.Content)
	}
	if split.PanelId != nil {
		q = q.Set("panel_id", split.PanelId)
	}
	if split.Effect != nil {
		q = q.Set("effect", split.Effect)
	}
	res, err := q.RunWith(i.db).Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) DeleteSplits(scriptId int64) error {
	_, err := sq.Delete("splits").
		Where(sq.Eq{"script_id": scriptId}).
		RunWith(i.db).
		Exec()
	return err
}

func (i *implementation) DeleteSplit(id int64) error {
	res, err := sq.Delete("splits").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetSplits(scriptId int64) ([]script.Split, error) {
	rows, err := sq.Select("id", "script_id", "content", "panel_id", "effect").
		From("splits").
		Where(sq.Eq{"script_id": scriptId}).
		RunWith(i.db).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	splits := []script.Split{}
	for rows.Next() {
		var s script.Split
		if err := rows.Scan(&s.Id, &s.ScriptId, &s.Content, &s.PanelId, &s.Effect); err != nil {
			return nil, err
		}
		splits = append(splits, s)
	}
	return splits, rows.Err()
}

func (i *implementation) GetSplit(id int64) (*script.Split, error) {
	var s script.Split
	err := sq.Select("id", "script_id", "content", "panel_id", "effect").
		From("splits").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		QueryRow().
		Scan(&s.Id, &s.ScriptId, &s.Content, &s.PanelId, &s.Effect)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
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
