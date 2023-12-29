package qb_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/t0k4r/qb"
)

func TestSelect(t *testing.T) {
	mfun := func(a any) any {
		switch a := a.(type) {
		case string:
			return strings.ReplaceAll(fmt.Sprint(a), "'", "''")
		default:
			return a
		}
	}
	_ = mfun
	q := qb.SelectNil("animes a").
		Cols("a.title").
		Wheref("a.title = '%v'", "cowboy").
		OrderBy("a.title asc").
		Limit("10").
		Sql()
	t.Log(q)
}

func TestInsert(t *testing.T) {
	q := qb.Insert("animes").
		Set("title", "cowboy").
		Setf("rating", "$1").
		Setf("lol", "(select id where yyz = '%v')", "notnot").
		Set("aired", time.Now()).
		Set("description", nil).
		OnConflict(qb.DoNothing).
		Sql()
	t.Log(q)
}

func TestDelete(t *testing.T) {
	q := qb.Delete("animes").
		Wheref("title = '%v'", "cowboy").
		Sql()
	t.Log(q)

}

func TestUpdate(t *testing.T) {
	q := qb.Update("animes").
		Set("title", "cowboy").
		Wheref("title = '%v'", "cowboy").
		Sql()
	t.Log(q)

}
