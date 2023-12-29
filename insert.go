package qb

import (
	"database/sql"
	"fmt"
	"strings"
)

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
	args       []any
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
func (q *QInsert) comma() {
	if q.cols.Len() != 0 {
		q.cols.WriteString(", ")
		q.vals.WriteString(", ")
	}
}

func (q *QInsert) Set(col string, val any) *QInsert {
	v := normalize(val)
	if v != "" {
		q.comma()
		q.cols.WriteString(col)
		q.vals.WriteString(v)
	}
	return q
}

func (q *QInsert) Setf(col string, format string, a ...any) *QInsert {
	q.comma()
	q.cols.WriteString(col)
	q.vals.WriteString(fmt.Sprintf(format, a...))
	return q
}

func (q *QInsert) Args(args ...any) *QInsert {
	q.args = append(q.args, args...)
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
	_, err := db.Exec(q.Sql(), append(q.args, args...)...)
	return err
}
