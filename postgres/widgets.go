package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bobisme/discoverable-api/models"
	"github.com/bobisme/discoverable-api/postgres/cursor"
	"github.com/bobisme/discoverable-api/storage"
)

//go:generate msgp

type WidgetCursor struct {
	CreatedAt int64  `msg:"created_at"`
	ID        string `msg:"id"`
}

func NewWidgetCursor(createdAt time.Time, id string) *WidgetCursor {
	return &WidgetCursor{
		CreatedAt: createdAt.UnixNano(),
		ID:        id,
	}
}

func NewWidgetCursorFromString(c string) *WidgetCursor {
	wc := new(WidgetCursor)
	err := cursor.Decode(c, wc)
	fatalIfErr(err)
	return wc
}

func (c *WidgetCursor) String() string {
	res, err := cursor.Encode(c)
	fatalIfErr(err)
	return res
}

type WidgetTable struct{}

func (t *WidgetTable) name() string { return "widgets" }

func (t *WidgetTable) columns() []string {
	return []string{
		"id",
		"name",
		"created_at",
	}
}

func (t *WidgetTable) cols() string { return strings.Join(t.columns(), ",") }

func (t *WidgetTable) tableCols() string {
	cols := make([]string, 0)
	for _, col := range t.columns() {
		cols = append(cols, t.name()+"."+col)
	}
	return strings.Join(cols, ",")
}

type WidgetStorage struct {
	table   WidgetTable
	db      *sql.DB
	counter *Counter
}

func (s *WidgetStorage) scan(rows *sql.Rows) ([]*models.Widget, error) {
	bound := []*models.Widget{}
	for rows.Next() {
		b := models.Widget{}
		if err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt); err != nil {
			return nil, err
		}
		bound = append(bound, &b)
	}
	return bound, nil
}

func (s *WidgetStorage) scanq(stmt string, args ...interface{}) ([]*models.Widget, error) {
	rows, err := s.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	return s.scan(rows)
}

func (s *WidgetStorage) Save(widget *models.Widget) (*models.Widget, error) {
	res, err := s.scanq(fmt.Sprintf(`
		INSERT INTO %[1]s (id,name,created_at)
		VALUES ($1, $2, NOW())
		RETURNING %[2]s`, s.table.name(), s.table.cols()),
		widget.ID, widget.Name)
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

func (s *WidgetStorage) Get(id string) (*models.Widget, error) {
	res, err := s.scanq(fmt.Sprintf(`
		SELECT %[1]s FROM %[2]s
		WHERE id = $1
		AND deleted_at IS NULL
		`, s.table.cols(), s.table.name()), id)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrNoRows
	} else if len(res) > 1 {
		return nil, ErrMultipleRows
	}
	return res[0], nil
}

func (s *WidgetStorage) List(cursor string, limit int) (storage.WidgetPage, error) {
	countStmt := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s WHERE deleted_at IS NULL
	`, s.table.name())
	queryStmt := fmt.Sprintf(`
		SELECT %[1]s FROM %[2]s
		LEFT JOIN widget_things
			ON widget_things.widget_id = id
		LEFT JOIN things
			ON things.id = widget_things.thing_id
		WHERE deleted_at IS NULL
		WHERE thing.deleted_at IS NULL
		ORDER BY created_at, id
		LIMIT $1`,
		s.table.cols(), s.table.name())
	res, err := s.scanq(queryStmt, limit)
	fatalIfErr(err)
	count, err := s.counter.Count(countStmt)
	fatalIfErr(err)
	return storage.WidgetPage{
		Items: res,
		Total: count,
	}, nil
}

func (s *WidgetStorage) Delete(id string) error {
	panic("not implemented")
}
