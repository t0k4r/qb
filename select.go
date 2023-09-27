package qb

import "database/sql"

type selectQuery struct {
}

func Select(table string) selectQuery {
	return selectQuery{}
}

func (s selectQuery) Join(table, on string) selectQuery {
	return s
}
func (s selectQuery) Column(column string) selectQuery {
	return s
}
func (s selectQuery) Where(where string) selectQuery {
	return s
}
func (s selectQuery) OrderBy(orderBy string) selectQuery {
	return s
}
func (s selectQuery) Limit(limit string) selectQuery {
	return s
}

func (s selectQuery) Sql() string {
	return ""
}

func (s selectQuery) Run(db *sql.DB, onRow func(*sql.Rows) error) error {
	rows, err := db.Query(s.Sql())
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
