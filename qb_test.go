package qb_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/t0k4r/qb"
)

type ds struct {
}

func (ds) Scan(*sql.Rows) (qb.Selectable, error) {
	return nil, nil
}

func TestSelect(t *testing.T) {
	q := qb.Select[ds]().
		Add(" select a.title from animes a ").
		Addf(" where a.title = '%v' ", "ok").
		Sql()
	t.Log(q)
}

func TestInsert(t *testing.T) {
	q := qb.Insert("animes").
		Add("title", "cowboy").
		Addf("rating", "$1").
		Add("aired", time.Now()).
		Add("description", nil).
		OnConflict(qb.DoNothing).
		Sql()
	t.Log(q)
}
