package qb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrDummyScan error = errors.New("ISelectable scanned dummy cant use Query and QueryFirst on this QSelect")

type dummy int

func (dummy) Scan(rows *sql.Rows) (ISelectable, error) {
	return dummy(0), ErrDummyScan
}

type ISelectable interface {
	Scan(*sql.Rows) (ISelectable, error)
}

type QSelect[T ISelectable] struct {
	from  string
	cols  strings.Builder
	join  strings.Builder
	where strings.Builder
	order strings.Builder
	limit strings.Builder
	args  []any
}

func SelectNil(from string) *QSelect[dummy] {
	return Select[dummy](from)
}

func Select[T ISelectable](from string) *QSelect[T] {
	return &QSelect[T]{from: from}
}

func (q *QSelect[T]) Cols(cols string) *QSelect[T] {
	q.cols.WriteString(" ")
	q.cols.WriteString(cols)
	q.cols.WriteString(" ")
	return q
}
func (q *QSelect[T]) Colsf(format string, a ...any) *QSelect[T] {
	q.cols.WriteString(" ")
	q.cols.WriteString(fmt.Sprintf(format, a...))
	q.cols.WriteString(" ")
	return q
}
func (q *QSelect[T]) Join(join string) *QSelect[T] {
	q.join.WriteString(" join ")
	q.join.WriteString(join)
	q.join.WriteString(" ")
	return q
}
func (q *QSelect[T]) Joinf(format string, a ...any) *QSelect[T] {
	q.join.WriteString(" join ")
	q.join.WriteString(fmt.Sprintf(format, a...))
	q.join.WriteString(" ")
	return q
}
func (q *QSelect[T]) LeftJoin(join string) *QSelect[T] {
	q.join.WriteString(" left join ")
	q.join.WriteString(join)
	q.join.WriteString(" ")
	return q
}
func (q *QSelect[T]) LeftJoinf(format string, a ...any) *QSelect[T] {
	q.join.WriteString(" left join ")
	q.join.WriteString(fmt.Sprintf(format, a...))
	q.join.WriteString(" ")
	return q
}
func (q *QSelect[T]) RightJoin(join string) *QSelect[T] {
	q.join.WriteString(" right join ")
	q.join.WriteString(join)
	q.join.WriteString(" ")
	return q
}
func (q *QSelect[T]) RightJoinf(format string, a ...any) *QSelect[T] {
	q.join.WriteString(" right join ")
	q.join.WriteString(fmt.Sprintf(format, a...))
	q.join.WriteString(" ")
	return q
}
func (q *QSelect[T]) Where(where string) *QSelect[T] {
	if q.where.Len() == 0 {
		q.where.WriteString(" where ")
	} else {
		q.where.WriteString(" ")
	}
	q.where.WriteString(where)
	q.where.WriteString(" ")
	return q
}
func (q *QSelect[T]) Wheref(format string, a ...any) *QSelect[T] {
	if q.where.Len() == 0 {
		q.where.WriteString(" where ")
	} else {
		q.where.WriteString(" ")
	}
	q.where.WriteString(fmt.Sprintf(format, a...))
	q.where.WriteString(" ")
	return q
}
func (q *QSelect[T]) OrderBy(oredr string) *QSelect[T] {
	if q.order.Len() == 0 {
		q.order.WriteString(" order by ")
	} else {
		q.order.WriteString(" ")
	}
	q.order.WriteString(oredr)
	q.order.WriteString(" ")
	return q
}
func (q *QSelect[T]) OrderByf(format string, a ...any) *QSelect[T] {
	if q.order.Len() == 0 {
		q.order.WriteString(" order by ")
	} else {
		q.order.WriteString(" ")
	}
	q.order.WriteString(fmt.Sprintf(format, a...))
	q.order.WriteString(" ")
	return q
}
func (q *QSelect[T]) Limit(limit string) *QSelect[T] {
	if q.limit.Len() == 0 {
		q.limit.WriteString(" limit ")
	} else {
		q.limit.WriteString(" ")
	}
	q.limit.WriteString(limit)
	q.limit.WriteString(" ")
	return q
}
func (q *QSelect[T]) Limitf(format string, a ...any) *QSelect[T] {
	if q.limit.Len() == 0 {
		q.limit.WriteString(" limit ")
	} else {
		q.limit.WriteString(" ")
	}
	q.limit.WriteString(fmt.Sprintf(format, a...))
	q.limit.WriteString(" ")
	return q
}

func (q *QSelect[T]) Args(args ...any) *QSelect[T] {
	q.args = append(q.args, args...)
	return q
}

func (q *QSelect[T]) Sql() string {
	query := strings.Builder{}
	query.WriteString("select")
	if q.cols.Len() == 0 {
		query.WriteString(" * ")
	} else {
		query.WriteString(q.cols.String())
	}
	query.WriteString("from ")
	query.WriteString(q.from)
	query.WriteString(q.join.String())
	query.WriteString(q.where.String())
	query.WriteString(q.order.String())
	query.WriteString(q.limit.String())
	return query.String()
}

func (q *QSelect[T]) Query(db *sql.DB, args ...any) ([]T, error) {
	var items []T
	rows, err := db.Query(q.Sql(), append(q.args, args...)...)
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
	rows, err := db.Query(q.Sql(), append(q.args, args...)...)
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
func (q *QSelect[T]) QueryEach(db *sql.DB, rowFunc func(*sql.Rows) error, args ...any) error {
	rows, err := db.Query(q.Sql(), append(q.args, args...)...)
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rowFunc(rows)
		if err != nil {
			return err
		}
	}
	return nil
}
