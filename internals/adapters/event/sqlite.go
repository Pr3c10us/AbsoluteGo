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
		Columns("status", "operation", "chapter_id", "script_id", "vab_id").
		Values(event.Status, event.Operation, event.ChapterId, event.ScriptId, event.VabId).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) UpdateStatus(id int64, status event.Status) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid status: %s", status)
	}

	q := sq.Update("events").Where(sq.Eq{"id": id}).Set("status", status)

	res, err := q.RunWith(i.db).Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetEvents(filter event.Filter) ([]event.Event, error) {
	q := sq.Select("id", "status", "operation", "chapter_id", "script_id", "vab_id").
		From("events")

	if filter.Status != "" {
		q = q.Where(sq.Eq{"status": filter.Status})
	}

	if filter.Operation != "" {
		q = q.Where(sq.Eq{"operation": filter.Operation})
	}

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
		if err = rows.Scan(&e.Id, &e.Status, &e.Operation, &e.ChapterId, &e.ScriptId, &e.VabId); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (i *implementation) GetEvent(id int64) (*event.Event, error) {
	var e event.Event
	err := sq.Select("id", "status", "operation", "chapter_id", "script_id", "vab_id").
		From("events").
		Where(sq.Eq{"id": id}).
		RunWith(i.db).
		QueryRow().
		Scan(&e.Id, &e.Status, &e.Operation, &e.ChapterId, &e.ScriptId, &e.VabId)

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
