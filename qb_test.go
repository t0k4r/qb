package qb_test

import (
	"testing"

	"github.com/t0k4r/qb"
)

func TestSelect(t *testing.T) {
	q := qb.Select("animes a").
		Cols("a.id", "a.title", "ai.aired").
		LJoin("anime_infos ai", "ai.anime_id=a.id").
		Where("a.title like ? and ai.aired > ?").
		OrderBy("ai.aired desc").
		Limit("10").
		Sql()
	t.Log(q)
}

func TestInsert(t *testing.T) {
	q := qb.Insert("animes").
		Col("title", "cowboy").
		Col("aired", qb.Arg("$1")).
		Col("aired_id",
			qb.Select("anime_infos ai").
				Cols("ai.aired_id").
				Where("ai.genre_id = $2")).
		Col("xd", "").
		Sql(qb.Replace)
	t.Log(q)

}
