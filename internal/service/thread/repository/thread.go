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

func jsonToDbModel(f *models.Thread) *Thread {
	return &Thread{
		Id:      f.Id,
		Title:   f.Title,
		Author:  f.Author,
		Forum:   f.Forum,
		Message: f.Message,
		Votes:   f.Votes,
		Slug:    f.Slug,
		Created: f.Created,
	}
}
func dbToJsonModel(f *Thread) *models.Thread {
	return &models.Thread{
		Id:      f.Id,
		Title:   f.Title,
		Author:  f.Author,
		Forum:   f.Forum,
		Message: f.Message,
		Votes:   f.Votes,
		Slug:    f.Slug,
		Created: f.Created,
	}
}
