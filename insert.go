package qb

import (
	"database/sql"
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

func (i QInsert) Build() BQInsert {
	return BQInsert{
		sql: i.Sql(),
	}
}

type BQInsert struct {
	sql string
}

func (i BQInsert) Exec(db *sql.DB, args ...any) error {
	_, err := db.Exec(i.sql, args...)
	return err
}
