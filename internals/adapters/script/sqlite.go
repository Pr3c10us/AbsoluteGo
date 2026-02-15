package script

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
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

func (i *implementation) GetScripts(query script.Query) ([]script.Script, error) {
	q := sq.Select("id", "name", "content", "book_id", "chapters").
		From("scripts")
	if query.BookId > 0 {
		q = q.Where(sq.Eq{"book_id": query.BookId})
	}
	if query.Name != "" {
		q = q.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", query.Name)})
	}
	if len(query.Ids) > 0 {
		q = q.Where(sq.Eq{"id": query.Ids})
	}
	if query.Chapter > 0 {
		q = q.Where(sq.Expr(
			"EXISTS (SELECT 1 FROM json_each(chapters) WHERE value = ?)",
			query.Chapter,
		))
	}

	rows, err := q.RunWith(i.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts []script.Script
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
	if errors.Is(err, sql.ErrNoRows) {
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
		Columns("script_id", "content", "previous_content", "panel_id", "effect").
		Values(split.ScriptId, split.Content, split.PreviousContent, split.PanelId, split.Effect).
		RunWith(i.db).
		Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (i *implementation) CreateManySplit(splits []script.Split) ([]script.Split, error) {
	if len(splits) == 0 {
		return nil, nil
	}

	tx, err := i.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := sq.Insert("splits").
		Columns("script_id", "content", "previous_content", "panel_id", "effect")

	for _, s := range splits {
		query = query.Values(s.ScriptId, s.Content, s.PreviousContent, s.PanelId, s.Effect)
	}

	res, err := query.RunWith(tx).Exec()
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// SQLite guarantees contiguous rowids for multi-row inserts.
	// LastInsertId() returns the ID of the *last* row inserted.
	firstID := lastID - int64(len(splits)) + 1

	result := make([]script.Split, len(splits))
	for idx, s := range splits {
		s.Id = firstID + int64(idx)
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
	if split.PreviousContent != nil {
		q = q.Set("previous_content", split.PreviousContent)
	}
	if split.PanelId != nil {
		q = q.Set("panel_id", split.PanelId)
	}
	if split.Effect != nil {
		q = q.Set("effect", split.Effect)
	}
	if split.AudioURL != nil {
		q = q.Set("audio_url", split.AudioURL)
	}
	if split.AudioDuration != nil {
		q = q.Set("audio_duration", split.AudioDuration)
	}
	if split.VideoURL != nil {
		q = q.Set("video_url", split.VideoURL)
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

func (i *implementation) DeleteSplit(ids []int64) error {
	res, err := sq.Delete("splits").
		Where(sq.Eq{"id": ids}).
		RunWith(i.db).
		Exec()
	if err != nil {
		return err
	}
	return assertRowAffected(res)
}

func (i *implementation) GetSplits(scriptId int64) ([]script.Split, error) {
	rows, err := sq.Select(
		"splits.id", "splits.script_id", "splits.content", "splits.previous_content", "splits.panel_id", "splits.effect", "splits.audio_url", "splits.audio_duration", "splits.video_url",
		"panels.id", "panels.page_id", "panels.url", "panels.panel_number", "panels.updated_at",
	).
		From("splits").
		LeftJoin("panels ON panels.id = splits.panel_id").
		Where(sq.Eq{"splits.script_id": scriptId}).
		RunWith(i.db).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splits []script.Split
	for rows.Next() {
		var s script.Split
		var panelId sql.NullInt64
		var panelPageId sql.NullInt64
		var panelURL sql.NullString
		var panelNumber sql.NullInt64
		var panelUpdatedAt sql.NullTime

		if err := rows.Scan(
			&s.Id, &s.ScriptId, &s.Content, &s.PreviousContent, &s.PanelId, &s.Effect, &s.AudioURL, &s.AudioDuration, &s.VideoURL,
			&panelId, &panelPageId, &panelURL, &panelNumber, &panelUpdatedAt,
		); err != nil {
			return nil, err
		}

		if panelId.Valid {
			url := &panelURL.String
			if !panelURL.Valid {
				url = nil
			}
			s.Panel = book.Panel{
				Id:          panelId.Int64,
				PageId:      panelPageId.Int64,
				URL:         url,
				PanelNumber: int(panelNumber.Int64),
				UpdatedAt:   panelUpdatedAt.Time,
			}
		}

		splits = append(splits, s)
	}
	return splits, rows.Err()
}

func (i *implementation) GetSplit(id int64) (*script.Split, error) {
	var s script.Split
	var panelId sql.NullInt64
	var panelPageId sql.NullInt64
	var panelURL sql.NullString
	var panelNumber sql.NullInt64
	var panelUpdatedAt sql.NullTime

	err := sq.Select(
		"splits.id", "splits.script_id", "splits.content", "splits.previous_content", "splits.panel_id", "splits.effect", "splits.audio_url", "splits.audio_duration", "splits.video_url",
		"panels.id", "panels.page_id", "panels.url", "panels.panel_number", "panels.updated_at",
	).
		From("splits").
		LeftJoin("panels ON panels.id = splits.panel_id").
		Where(sq.Eq{"splits.id": id}).
		RunWith(i.db).
		QueryRow().
		Scan(
			&s.Id, &s.ScriptId, &s.Content, &s.PreviousContent, &s.PanelId, &s.Effect, &s.AudioURL, &s.AudioDuration, &s.VideoURL,
			&panelId, &panelPageId, &panelURL, &panelNumber, &panelUpdatedAt,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if panelId.Valid {
		url := &panelURL.String
		if !panelURL.Valid {
			url = nil
		}
		s.Panel = book.Panel{
			Id:          panelId.Int64,
			PageId:      panelPageId.Int64,
			URL:         url,
			PanelNumber: int(panelNumber.Int64),
			UpdatedAt:   panelUpdatedAt.Time,
		}
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
