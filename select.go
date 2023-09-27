package qb

import (
	"database/sql"
	"strings"
)

type QSelect struct {
	table   string
	col     strings.Builder
	join    strings.Builder
	where   string
	limit   string
	orderBy string
	args    []any
}

func Select(table string) QSelect {
	return QSelect{table: table}
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
	s.where = where
	return s
}
func (s QSelect) OrderBy(orderBy string) QSelect {
	s.orderBy = orderBy
	return s
}
func (s QSelect) Limit(limit string) QSelect {
	s.limit = limit
	return s
}
func (s QSelect) Args(args ...any) QSelect {
	s.args = args
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
	query.WriteString("where ")
	query.WriteString(s.where)
	query.WriteString(" \norder by ")
	query.WriteString(s.orderBy)
	query.WriteString(" \nlimit ")
	query.WriteString(s.limit)
	query.WriteString(" \n")
	return query.String()
}

func (s QSelect) Build() BQSelect {
	return BQSelect{
		sql:  s.Sql(),
		args: s.args,
	}
}

type BQSelect struct {
	sql  string
	args []any
}

func (s BQSelect) Query(db *sql.DB, onRow func(*sql.Rows) error) error {
	rows, err := db.Query(s.sql, s.args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = onRow(rows)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s BQSelect) QueryRow(db *sql.DB, onRow func(*sql.Row) error) error {
	row := db.QueryRow(s.sql, s.args...)
	if row.Err() != nil {
		return row.Err()
	}
	return onRow(row)
}

func tst() []int {
	var ints []int
	Select("animes a").
		Cols("a.id", "a.title", "ai.aired").
		Join("anime_infos ai", "ai.anime_id=a.id").
		Where("a.title like ? and a.aired > ?").
		OrderBy("ai.aired desc").
		Limit("10").
		Args("%Cowboy%", "2023-14-11").
		Build().
		Query(nil, func(r *sql.Rows) error {
			var i int
			err := r.Scan(&i)
			if err != nil {
				return err
			}
			ints = append(ints, i)
			return nil
		})
	return ints

}
