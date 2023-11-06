package qb

import (
	"database/sql"
	"fmt"
	"strings"
)

type Arg string

type Conflict string

const (
	Ignore  Conflict = "ignore"
	Replace Conflict = "replace"
	none    Conflict = "none"
)

type QInsert struct {
	table string
	cols  strings.Builder
	vals  strings.Builder
}

func Insert(into string) *QInsert {
	return &QInsert{
		table: into,
		cols:  strings.Builder{},
		vals:  strings.Builder{},
	}
}

func (i *QInsert) insComma() {
	if i.cols.Len() != 0 {
		i.cols.WriteString(", ")
		i.vals.WriteString(", ")
	}
}

func (i *QInsert) Col(col string, value any) *QInsert {
	switch value := value.(type) {
	case *int, *uint, *float32, *float64:
		if value != nil {
			i.insComma()
			i.cols.WriteString(col)
			i.vals.WriteString(fmt.Sprint(value))
		}
	case int, uint, float32, float64:
		i.insComma()
		i.cols.WriteString(col)
		i.vals.WriteString(fmt.Sprint(value))
	case *string:
		if value != nil && *value != "" {
			i.insComma()
			i.cols.WriteString(col)
			i.vals.WriteString("'")
			i.vals.WriteString(strings.ReplaceAll(*value, "'", "''"))
			i.vals.WriteString("'")
		}
	case string:
		if value != "" {
			i.insComma()
			i.cols.WriteString(col)
			i.vals.WriteString("'")
			i.vals.WriteString(strings.ReplaceAll(value, "'", "''"))
			i.vals.WriteString("'")
		}
	case Arg:
		i.insComma()
		i.cols.WriteString(col)
		i.vals.WriteString(string(value))

	case *QSelect:
		i.insComma()
		i.cols.WriteString(col)
		i.vals.WriteString("(")
		i.vals.WriteString(value.Sql())
		i.vals.WriteString(")")
	default:
		panic(value)
	}
	return i
}

func (i *QInsert) Sql(mod ...Conflict) string {
	var query strings.Builder
	var fmod Conflict
	if len(mod) != 0 {
		fmod = mod[0]
	} else {
		fmod = none
	}
	query.WriteString("insert into ")
	query.WriteString(i.table)
	query.WriteString("(")
	query.WriteString(i.cols.String())
	query.WriteString(") values ")
	query.WriteString("(")
	query.WriteString(i.vals.String())
	switch fmod {
	case Ignore:
		query.WriteString(") on conflict do nothing ")
	case Replace:
		query.WriteString(") on conflict do update set ")
		strip := strings.ReplaceAll(i.cols.String(), " ", "")
		cols := strings.Split(strip, ",")
		for i, col := range cols {
			if i != 0 {
				query.WriteString(", ")
			}
			query.WriteString(fmt.Sprintf("%v=EXCLUDED.%v", col, col))
			_ = col
		}
	case none:
		query.WriteString(")")
	}
	return query.String()
}

type Insertable interface {
	Dump() (cols []string, vals []any)
}

func Exec[T Insertable](item T, into string, db *sql.DB) error {
	q := Insert(into)
	cols, vals := item.Dump()
	for i, col := range cols {
		val := vals[i]
		q.Col(col, val)
	}
	_, err := db.Exec(q.Sql())
	return err
}

type QSelect struct {
	table   string
	col     strings.Builder
	join    strings.Builder
	where   strings.Builder
	limit   strings.Builder
	orderBy strings.Builder
}

func Select(from string) *QSelect {
	return &QSelect{table: from}
}

func (s *QSelect) LJoin(table, on string) *QSelect {
	s.join.WriteString(" left join ")
	s.join.WriteString(table)
	s.join.WriteString(" on ")
	s.join.WriteString(on)
	return s
}

func (s *QSelect) RJoin(table, on string) *QSelect {
	s.join.WriteString(" right join ")
	s.join.WriteString(table)
	s.join.WriteString(" on ")
	s.join.WriteString(on)
	return s
}

func (s *QSelect) Join(table, on string) *QSelect {
	s.join.WriteString(" join ")
	s.join.WriteString(table)
	s.join.WriteString(" on ")
	s.join.WriteString(on)
	return s
}
func (s *QSelect) Cols(cols ...string) *QSelect {
	for i, col := range cols {
		if i != 0 {
			s.col.WriteString(", ")
		}
		s.col.WriteString(col)
	}
	return s
}
func (s *QSelect) Where(where string) *QSelect {
	s.where.WriteString(where)
	return s
}

func (s *QSelect) Wheref(wherefmt string, args ...any) *QSelect {
	for i, a := range args {
		switch a := a.(type) {
		case string:
			args[i] = strings.ReplaceAll(a, "'", "''")
		}
	}
	s.where.WriteString(fmt.Sprintf(wherefmt, args...))
	return s
}
func (s *QSelect) OrderBy(orderBy string) *QSelect {
	s.orderBy.WriteString(orderBy)
	return s
}
func (s *QSelect) Limit(limit string) *QSelect {
	s.limit.WriteString(limit)
	return s
}
func (s *QSelect) Sql() string {
	var query strings.Builder
	query.WriteString("select ")
	query.WriteString(s.col.String())
	query.WriteString(" from ")
	query.WriteString(s.table)
	query.WriteString(s.join.String())
	if s.where.Len() != 0 {
		query.WriteString(" where ")
		query.WriteString(s.where.String())
	}
	if s.orderBy.Len() != 0 {
		query.WriteString(" order by ")
		query.WriteString(s.orderBy.String())
	}
	if s.limit.Len() != 0 {
		query.WriteString(" limit ")
		query.WriteString(s.limit.String())
	}
	return query.String()
}

type Selectable interface {
	Scan(*sql.Rows) (Selectable, error)
}

func Query[T Selectable](query *QSelect, db *sql.DB, args ...any) ([]T, error) {
	var items []T
	rows, err := db.Query(query.Sql(), args...)
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
