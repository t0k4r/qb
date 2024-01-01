package qb

import (
	"database/sql"
	"fmt"
	"strings"
)

type QUpdate struct {
	table string
	set   strings.Builder
	where strings.Builder
	args  []any
}

func Update(table string) *QUpdate {
	return &QUpdate{table: table}
}

func (q *QUpdate) comma() {
	if q.set.Len() != 0 {
		q.set.WriteString(",")
	}
}

func (q *QUpdate) Set(col string, val string) *QUpdate {
	q.comma()
	q.set.WriteString(" ")
	q.set.WriteString(col)
	q.set.WriteString("=")
	q.set.WriteString(val)
	q.set.WriteString(" ")
	return q
}

func (q *QUpdate) Setn(col string, val any) *QUpdate {
	return q.Set(col, normalize(val))
}

func (q *QUpdate) Setf(col string, format string, a ...any) *QUpdate {
	return q.Set(col, fmt.Sprintf(format, a...))
}
func (q *QUpdate) Seta(col string, val string, args ...any) *QUpdate {
	return q.Set(col, val).Args(args...)
}

func (q *QUpdate) Where(where string) *QUpdate {
	if q.where.Len() == 0 {
		q.where.WriteString(" where ")
	} else {
		q.where.WriteString(" ")
	}
	q.where.WriteString(where)
	q.where.WriteString(" ")
	return q
}

func (q *QUpdate) Wheref(format string, a ...any) *QUpdate {
	return q.Where(fmt.Sprintf(format, a...))
}
func (q *QUpdate) Wherea(where string, args ...any) *QUpdate {
	return q.Where(where).Args(args...)
}

func (q *QUpdate) Args(args ...any) *QUpdate {
	q.args = append(q.args, args...)
	return q
}

func (q *QUpdate) Sql() string {
	query := strings.Builder{}
	query.WriteString("update ")
	query.WriteString(q.table)
	query.WriteString(" set ")
	query.WriteString(q.set.String())
	query.WriteString(q.where.String())
	return query.String()
}

func (q *QUpdate) Exec(db *sql.DB, args ...any) error {
	_, err := db.Exec(q.Sql(), append(q.args, args...)...)
	return err
}
