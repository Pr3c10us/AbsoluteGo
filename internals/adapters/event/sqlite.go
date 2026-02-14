package event

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/event"

	sq "github.com/Masterminds/squirrel"
)

type implementation struct {
	db *sql.DB
}

func NewEventImplementation(db *sql.DB) event.Interface {
	return &implementation{
		db: db,
	}
}

func (i *implementation) Create(event event.Event) (int64, error) {
	res, err := sq.Insert("events").
		Columns("status", "operation", "description", "book_id", "chapter_id", "script_id", "vab_id").
		Values(event.Status, event.Operation, event.Description, event.BookId, event.ChapterId, event.ScriptId, event.VabId).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) Update(id int64, e event.Event) error {
	if e.Status != "" && !e.Status.IsValid() {
		return fmt.Errorf("invalid status: %s", e.Status)
	}

	q := sq.Update("events").Where(sq.Eq{"id": id})

	if e.Status != "" {
		q = q.Set("status", e.Status)
	}
	if e.Operation != "" {
		q = q.Set("operation", e.Operation)
	}

	if e.Description != "" {
		q = q.Set("description", e.Description)
	}
	if e.BookId != 0 {
		q = q.Set("book_id", e.BookId)
	}
	if e.ChapterId != 0 {
		q = q.Set("chapter_id", e.ChapterId)
	}
	if e.ScriptId != 0 {
		q = q.Set("script_id", e.ScriptId)
	}
	if e.VabId != 0 {
		q = q.Set("vab_id", e.VabId)
	}

	res, err := q.RunWith(i.db).Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetEvents(filter event.Filter) ([]event.Event, error) {
	q := sq.Select("id", "status", "operation", "description", "book_id", "chapter_id", "script_id", "vab_id", "updated_at").
		From("events")

	if filter.Status != "" {
		q = q.Where(sq.Eq{"status": filter.Status})
	}

	if filter.Operation != "" {
		q = q.Where(sq.Eq{"operation": filter.Operation})
	}

	q = q.OrderBy("updated_at DESC")

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}

	q = q.Limit(uint64(limit))
	offset := (page - 1) * limit
	if offset > 0 {
		q = q.Offset(uint64(offset))
	}

	rows, err := q.RunWith(i.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var e event.Event
		if err = rows.Scan(&e.Id, &e.Status, &e.Operation, &e.Description, &e.BookId, &e.ChapterId, &e.ScriptId, &e.VabId, &e.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (i *implementation) GetEvent(id int64) (*event.Event, error) {
	var e event.Event
	err := sq.Select("id", "status", "operation", "description", "book_id", "chapter_id", "script_id", "vab_id", "updated_at").
		From("events").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		QueryRow().
		Scan(&e.Id, &e.Status, &e.Operation, &e.Description, &e.BookId, &e.ChapterId, &e.ScriptId, &e.VabId, &e.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &e, nil
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
