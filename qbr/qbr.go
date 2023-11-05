package qbr

import (
	"database/sql"

	"github.com/t0k4r/qb"
)

type Selectable interface {
	Cols() qb.QSelect
	Scan(*sql.Rows) (Selectable, error)
}

type QSelect[T Selectable] struct {
	qb.QSelect
}

func Select[T Selectable]() QSelect[T] {
	var t T
	return QSelect[T]{QSelect: t.Cols()}
}

func (s QSelect[T]) Where(where string) QSelect[T] {
	s.QSelect = s.QSelect.Where(where)
	return s
}

func (s QSelect[T]) Wheref(wherefmt string, args ...any) QSelect[T] {
	s.QSelect = s.QSelect.Wheref(wherefmt, args...)
	return s
}
func (s QSelect[T]) OrderBy(orderBy string) QSelect[T] {
	s.QSelect = s.QSelect.OrderBy(orderBy)
	return s
}
func (s QSelect[T]) Limit(limit string) QSelect[T] {
	s.QSelect = s.QSelect.Limit(limit)
	return s
}

func (s QSelect[T]) Query(db *sql.DB, args ...any) ([]T, error) {
	var items []T
	rows, err := db.Query(s.Sql(), args...)
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

type Col struct {
	Column string
	Value  any
}

type Insertable interface {
	Cols() []Col
}

type QInsert struct {
	qb.QInsert
}

func Insert[T Insertable](into string, item T) QInsert {
	ins := qb.Insert(into)
	for _, i := range item.Cols() {
		ins = ins.Col(i.Column, i.Value)
	}
	return QInsert{QInsert: ins}
}

func (i QInsert) Exec(db *sql.DB, args ...any) error {
	_, err := db.Exec(i.Sql(), args...)
	return err
}
