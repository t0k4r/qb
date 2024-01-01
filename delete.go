package qb

import (
	"database/sql"
	"fmt"
	"strings"
)

type QDelete struct {
	from  string
	where strings.Builder
	args  []any
}

func Delete(from string) *QDelete {
	return &QDelete{from: from}
}

func (q *QDelete) Where(where string) *QDelete {
	if q.where.Len() == 0 {
		q.where.WriteString(" where ")
	} else {
		q.where.WriteString(" ")
	}
	q.where.WriteString(where)
	q.where.WriteString(" ")
	return q
}

func (q *QDelete) Wheref(format string, a ...any) *QDelete {
	return q.Where(fmt.Sprintf(format, a...))
}

func (q *QDelete) Wherea(where string, args ...any) *QDelete {
	return q.Where(where).Args(args...)
}

func (q *QDelete) Args(args ...any) *QDelete {
	q.args = append(q.args, args...)
	return q
}

func (q *QDelete) Sql() string {
	query := strings.Builder{}
	query.WriteString("delete from ")
	query.WriteString(q.from)
	query.WriteString(q.where.String())
	return query.String()
}

func (q *QDelete) Exec(db *sql.DB, args ...any) error {
	_, err := db.Exec(q.Sql(), append(q.args, args...)...)
	return err
}
