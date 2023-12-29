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

func (q *QUpdate) Set(col string, val any) *QUpdate {
	v := normalize(val)
	if v == "" {
		v = "NULL"
	}
	q.comma()
	q.set.WriteString(" ")
	q.set.WriteString(col)
	q.set.WriteString("=")
	q.set.WriteString(v)
	q.set.WriteString(" ")
	return q
}

func (q *QUpdate) Setf(col string, format string, a ...any) *QUpdate {
	q.comma()
	q.set.WriteString(" ")
	q.set.WriteString(col)
	q.set.WriteString("=")
	q.set.WriteString(fmt.Sprintf(format, a...))
	q.set.WriteString(" ")
	return q
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
	if q.where.Len() == 0 {
		q.where.WriteString(" where ")
	} else {
		q.where.WriteString(" ")
	}
	q.where.WriteString(fmt.Sprintf(format, a...))
	q.where.WriteString(" ")
	return q
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
