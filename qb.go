package qb

import (
	"fmt"
	"strings"
)

type Arg string

type QInsert struct {
	table string
	cols  strings.Builder
	vals  strings.Builder
}

func Insert(into string) QInsert {
	return QInsert{table: into}
}

func (i *QInsert) comma() {
	if i.cols.Len() != 0 {
		i.cols.WriteString(", ")
		i.vals.WriteString(", ")
	}
}

func (i QInsert) Col(col string, value any) QInsert {
	switch value := value.(type) {
	case int, uint, float32, float64:
		i.comma()
		i.cols.WriteString(col)
		i.vals.WriteString(fmt.Sprint(value))
	case string:
		if value != "" {
			i.comma()
			i.cols.WriteString(col)
			i.vals.WriteString("'")
			i.vals.WriteString(strings.ReplaceAll(value, "'", "''"))
			i.vals.WriteString("'")
		}
	case Arg:
		i.comma()
		i.cols.WriteString(col)
		i.vals.WriteString(string(value))

	case QSelect:
		i.comma()
		i.cols.WriteString(col)
		i.vals.WriteString("(\n")
		i.vals.WriteString(value.Sql())
		i.vals.WriteString(")")
	default:
		panic(value)
	}
	return i
}

func (i QInsert) Sql() string {
	var query strings.Builder
	query.WriteString("insert into ")
	query.WriteString(i.table)
	query.WriteString("(")
	query.WriteString(i.cols.String())
	query.WriteString(")\nvalues ")
	query.WriteString("(")
	query.WriteString(i.vals.String())
	query.WriteString(") on conflict do nothing\n")
	return query.String()
}

type QSelect struct {
	table   string
	col     strings.Builder
	join    strings.Builder
	where   strings.Builder
	limit   strings.Builder
	orderBy strings.Builder
}

func Select(from string) QSelect {
	return QSelect{table: from}
}

func (s QSelect) LJoin(table, on string) QSelect {
	s.join.WriteString("left join ")
	s.join.WriteString(table)
	s.join.WriteString(" on ")
	s.join.WriteString(on)
	s.join.WriteString(" \n")
	return s
}

func (s QSelect) RJoin(table, on string) QSelect {
	s.join.WriteString("right join ")
	s.join.WriteString(table)
	s.join.WriteString(" on ")
	s.join.WriteString(on)
	s.join.WriteString(" \n")
	return s
}

func (s QSelect) Join(table, on string) QSelect {
	s.join.WriteString("join ")
	s.join.WriteString(table)
	s.join.WriteString(" on ")
	s.join.WriteString(on)
	s.join.WriteString(" \n")
	return s
}
func (s QSelect) Cols(cols ...string) QSelect {
	for i, col := range cols {
		if i != 0 {
			s.col.WriteString(", ")
		}
		s.col.WriteString(col)
	}
	return s
}
func (s QSelect) Where(where string) QSelect {
	s.where.WriteString(where)
	return s
}
func (s QSelect) OrderBy(orderBy string) QSelect {
	s.orderBy.WriteString(orderBy)
	return s
}
func (s QSelect) Limit(limit string) QSelect {
	s.limit.WriteString(limit)
	return s
}
func (s QSelect) Sql() string {
	var query strings.Builder
	query.WriteString("select ")
	query.WriteString(s.col.String())
	query.WriteString(" from ")
	query.WriteString(s.table)
	query.WriteString(" \n")
	query.WriteString(s.join.String())
	if s.where.Len() != 0 {
		query.WriteString("where ")
		query.WriteString(s.where.String())
	}
	if s.orderBy.Len() != 0 {
		query.WriteString(" \norder by ")
		query.WriteString(s.orderBy.String())
	}
	if s.limit.Len() != 0 {
		query.WriteString(" \nlimit ")
		query.WriteString(s.limit.String())
	}
	query.WriteString(" \n")
	return query.String()
}
