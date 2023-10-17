package qb

import (
	"fmt"
	"strings"
)

type Arg string

type Conflict string

const (
	Ignore  Conflict = "ignore"
	Replace Conflict = "replace"
)

type Dialect string

const (
	Postgres Dialect = "postgres"
	Sqlite   Dialect = "sqlite"
)

var DefaultDialect = Postgres

type QInsert struct {
	table   string
	cols    strings.Builder
	vals    strings.Builder
	Dialect Dialect
}

func Insert(into string) QInsert {
	return QInsert{
		table:   into,
		cols:    strings.Builder{},
		vals:    strings.Builder{},
		Dialect: DefaultDialect,
	}
}

func (i *QInsert) insComma() {
	if i.cols.Len() != 0 {
		i.cols.WriteString(", ")
		i.vals.WriteString(", ")
	}
}

func (i QInsert) Col(col string, value any) QInsert {
	switch value := value.(type) {
	case int, uint, float32, float64:
		i.insComma()
		i.cols.WriteString(col)
		i.vals.WriteString(fmt.Sprint(value))
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

	case QSelect:
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

func (i QInsert) middleSql(query *strings.Builder) {
	query.WriteString(i.table)
	query.WriteString("(")
	query.WriteString(i.cols.String())
	query.WriteString(") values ")
	query.WriteString("(")
	query.WriteString(i.vals.String())
}

func (i QInsert) Sql(mod ...Conflict) string {
	var query strings.Builder
	if len(mod) != 0 {
		mod := mod[0]
		switch i.Dialect {
		case Postgres:
			query.WriteString("insert into ")
			i.middleSql(&query)
			switch mod {
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
			}
			return query.String()

		case Sqlite:
			switch mod {
			case Ignore:
				query.WriteString("insert or ignore into ")
			case Replace:
				query.WriteString("insert or replace into ")
			}
			i.middleSql(&query)
			query.WriteString(")")
			return query.String()

		}
	}
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
	s.join.WriteString(" left join ")
	s.join.WriteString(table)
	s.join.WriteString(" on ")
	s.join.WriteString(on)
	return s
}

func (s QSelect) RJoin(table, on string) QSelect {
	s.join.WriteString(" right join ")
	s.join.WriteString(table)
	s.join.WriteString(" on ")
	s.join.WriteString(on)
	return s
}

func (s QSelect) Join(table, on string) QSelect {
	s.join.WriteString(" join ")
	s.join.WriteString(table)
	s.join.WriteString(" on ")
	s.join.WriteString(on)
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

func (s QSelect) Wheref(wherefmt string, args ...any) QSelect {
	s.where.WriteString(fmt.Sprintf(wherefmt, args...))
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
