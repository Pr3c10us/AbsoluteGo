package vab

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Pr3c10us/absolutego/internals/domains/vab"
)

type implementation struct {
	db *sql.DB
}

func NewVABImplementation(db *sql.DB) vab.Interface {
	return &implementation{db: db}
}

func (i *implementation) Create(v vab.VAB) (int64, error) {
	res, err := sq.Insert("vabs").
		Columns("name", "url", "script_id", "book_id").
		Values(v.Name, v.Url, v.ScriptId, v.BookId).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) Update(id int64, v vab.VAB) error {
	q := sq.Update("vabs").Where(sq.Eq{"id": id})
	if v.Name != "" {
		q = q.Set("name", v.Name)
	}
	if v.Url != nil {
		q = q.Set("url", v.Url)
	}
	if v.ScriptId > 0 {
		q = q.Set("script_id", v.ScriptId)
	}
	if v.BookId > 0 {
		q = q.Set("book_id", v.BookId)
	}
	res, err := q.RunWith(i.db).Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetVABs(name string, scriptId int64, bookId int64) ([]vab.VAB, error) {
	q := sq.Select("id", "name", "url", "script_id", "book_id").From("vabs")
	if name != "" {
		q = q.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", name)})
	}
	if scriptId > 0 {
		q = q.Where(sq.Eq{"script_id": scriptId})
	}
	if bookId > 0 {
		q = q.Where(sq.Eq{"book_id": bookId})
	}

	rows, err := q.RunWith(i.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vabs []vab.VAB
	for rows.Next() {
		var v vab.VAB
		if err := rows.Scan(&v.Id, &v.Name, &v.Url, &v.ScriptId, &v.BookId); err != nil {
			return nil, err
		}
		vabs = append(vabs, v)
	}
	return vabs, rows.Err()
}

func (i *implementation) Delete(id int64) error {
	res, err := sq.Delete("vabs").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) DeleteByScript(scriptId int64) error {
	res, err := sq.Delete("vabs").
		Where(sq.Eq{"script_id": scriptId}).
		RunWith(i.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetVAB(id int64) (*vab.VAB, error) {
	var v vab.VAB
	err := sq.Select("id", "name", "url", "script_id", "book_id").
		From("vabs").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		QueryRow().
		Scan(&v.Id, &v.Name, &v.Url, &v.ScriptId, &v.BookId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
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
