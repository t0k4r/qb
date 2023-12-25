package qb

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func normalize(val any) string {
	switch val := val.(type) {
	case nil:
		return "NULL"
	case []byte:
		return fmt.Sprintf("'%s'", val)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprint(val)
	case string:
		if val != "" {
			return fmt.Sprintf("'%v'", strings.ReplaceAll(val, "'", "''"))
		}
	case time.Time:
		return fmt.Sprintf("'%v'", val.Format("2006-01-02 15:04:05"))
	case *int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64,
		*float32, *float64:
		if val != nil {
			return fmt.Sprint(val)
		}
	case *string:
		if val != nil && *val != "" {
			return fmt.Sprintf("'%v'", strings.ReplaceAll(*val, "'", "''"))
		}
	case *time.Time:
		if val != nil {
			return fmt.Sprintf("'%v'", val.Format("2006-01-02 15:04:05"))
		}
	default:
		panic("unknown val type")
	}
	return ""
}

type Selectable interface {
	Scan(*sql.Rows) (Selectable, error)
}

type QSelect[T Selectable] struct {
	sql strings.Builder
}

func Select[T Selectable]() *QSelect[T] {
	return new(QSelect[T])
}

func (q *QSelect[T]) Add(str string) *QSelect[T] {
	q.sql.WriteString(" ")
	q.sql.WriteString(str)
	q.sql.WriteString(" ")
	return q
}

// add fmt
func (q *QSelect[T]) Addf(format string, a ...any) *QSelect[T] {
	q.sql.WriteString(" ")
	q.sql.WriteString(fmt.Sprintf(format, a...))
	q.sql.WriteString(" ")
	return q
}

// add normalize fmt
func (q *QSelect[T]) Addn(format string, a ...any) *QSelect[T] {
	q.sql.WriteString(" ")
	var args []any
	for _, a := range a {
		args = append(args, normalize(a))
	}
	q.Addf(format, args...)
	q.sql.WriteString(" ")
	return q
}

func (q *QSelect[T]) Sql() string {
	return q.sql.String()
}

func (q *QSelect[T]) Query(db *sql.DB, args ...any) ([]T, error) {
	var items []T
	rows, err := db.Query(q.Sql(), args...)
	if err != nil {
		return items, err
	}
	var i T
	for rows.Next() {
		item, err := i.Scan(rows)
		if err != nil {
			return items, err
		}
		items = append(items, item.(T))
	}
	return items, nil
}
func (q *QSelect[T]) QueryFirst(db *sql.DB, args ...any) (*T, error) {
	rows, err := db.Query(q.Sql(), args...)
	if err != nil {
		return nil, err
	}
	var item T
	if rows.Next() {
		s, err := item.Scan(rows)
		if err != nil {
			return nil, nil
		}
		item = s.(T)
		return &item, nil
	}
	return nil, nil
}

type OnConflict string

const (
	DoNothing OnConflict = "nothing"
	DoUpdate  OnConflict = "doUpdate"
	None      OnConflict = "none"
)

type QInsert struct {
	into       string
	onConflict OnConflict
	cols       strings.Builder
	vals       strings.Builder
}

func Insert(into string) *QInsert {
	return &QInsert{
		into:       into,
		onConflict: None,
		cols:       strings.Builder{},
		vals:       strings.Builder{},
	}
}
func (q *QInsert) OnConflict(do OnConflict) *QInsert {
	q.onConflict = do
	return q
}
func (q *QInsert) write(col string, val string) {
	if q.cols.Len() == 0 {
		q.cols.WriteString(col)
		q.vals.WriteString(val)
	} else {
		q.cols.WriteString(", " + col)
		q.vals.WriteString(", " + val)
	}
}

func (q *QInsert) Add(col string, val any) *QInsert {
	v := normalize(val)
	if v != "" {
		q.write(col, v)
	}
	return q
}

// add fmt
func (q *QInsert) Addf(col string, format string, a ...any) *QInsert {
	q.write(col, fmt.Sprintf(format, a...))
	return q
}

// add normalize fmt
func (q *QInsert) Addn(col string, format string, a ...any) *QInsert {
	var args []any
	for _, a := range a {
		args = append(args, normalize(a))
	}
	return q.Addf(col, format, args...)
}

// add raw
func (q *QInsert) Addr(col string, val any) *QInsert {
	q.write(col, fmt.Sprint(val))
	return q
}

func (q *QInsert) Sql() string {
	query := strings.Builder{}
	query.WriteString("insert into ")
	query.WriteString(q.into)
	query.WriteString("(")
	query.WriteString(q.cols.String())
	query.WriteString(") values ")
	query.WriteString("(")
	query.WriteString(q.vals.String())
	switch q.onConflict {
	case DoNothing:
		query.WriteString(") on conflict do nothing ")
	case DoUpdate:
		query.WriteString(") on conflict do update set ")
		strip := strings.ReplaceAll(q.cols.String(), " ", "")
		cols := strings.Split(strip, ",")
		for i, col := range cols {
			if i != 0 {
				query.WriteString(", ")
			}
			query.WriteString(fmt.Sprintf("%v=EXCLUDED.%v", col, col))
		}
	case None:
		query.WriteString(")")
	default:
		panic("unknown on conflict")
	}
	return query.String()
}

func (q *QInsert) Exec(db *sql.DB, args ...any) error {
	_, err := db.Exec(q.Sql(), args...)
	return err
}
