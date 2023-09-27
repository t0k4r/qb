package qb_test

import (
	"qb"
	"testing"
)

func TestSelect(t *testing.T) {

	q := qb.Select("animes a").
		Cols("a.id", "a.title", "ai.aired").
		Join("anime_infos ai", "ai.anime_id=a.id").
		Where("a.title like ? and a.aired > ?").
		OrderBy("ai.aired desc").
		Limit("10").Sql()

	t.Log(q)
}
