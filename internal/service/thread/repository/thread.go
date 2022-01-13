package repository

import (
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"time"
)

type Thread struct {
	Id      int32     `db:"id"`
	Title   string    `db:"title"`
	Author  string    `db:"author"`
	Forum   string    `db:"forum"`
	Message string    `db:"message"`
	Votes   int32     `db:"votes"`
	Slug    string    `db:"slug"`
	Created time.Time `db:"created"`
}

func jsonToDbModel(t *models.Thread) *Thread {
	return &Thread{
		Id:      t.Id,
		Title:   t.Title,
		Author:  t.Author,
		Forum:   t.Forum,
		Message: t.Message,
		Votes:   t.Votes,
		Slug:    t.Slug,
		Created: t.Created,
	}
}
func dbToJsonModel(t *Thread) *models.Thread {
	return &models.Thread{
		Id:      t.Id,
		Title:   t.Title,
		Author:  t.Author,
		Forum:   t.Forum,
		Message: t.Message,
		Votes:   t.Votes,
		Slug:    t.Slug,
		Created: t.Created,
	}
}
