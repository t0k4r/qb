package qb_test

import (
	"database/sql"
	"testing"

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
	// q := qb.Select("animes a").
	// 	Cols("a.id", "a.title", "ai.aired").
	// 	LJoin("anime_infos ai", "ai.anime_id=a.id").
	// 	Where("a.title like ? and ai.aired > ?").
	// 	OrderBy("ai.aired desc").
	// 	Limit("10").
	// 	Sql()
	// t.Log(q)
}

func TestInsert(t *testing.T) {
	q := qb.Insert("animes").
		OnConflict(qb.DoUpdate).
		Add("title", "cowboy").
		Addf("aired", "$1").
		Sql()
	t.Log(q)
	// q := qb.Insert("animes").
	// 	Col("title", "cowboy").
	// 	Col("aired", qb.Arg("$1")).
	// 	Col("aired_id",
	// 		qb.Select("anime_infos ai").
	// 			Cols("ai.aired_id").
	// 			Where("ai.genre_id = $2")).
	// 	Col("xd", "").
	// 	Sql(qb.Replace)
	// t.Log(q)

}
