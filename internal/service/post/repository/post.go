package repository

import (
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"time"
)

type Post struct {
	Id       int64     `db:"id"`
	Parent   int64     `db:"parent"`
	Author   string    `db:"author"`
	Message  string    `db:"message"`
	IsEdited bool      `db:"isedited"`
	Forum    string    `db:"forum"`
	Thread   int32     `db:"thread"`
	Created  time.Time `db:"created"`
}

func JsonToDbModel(p *models.Post) *Post {
	return &Post{
		Id:       p.Id,
		Parent:   p.Parent,
		Author:   p.Author,
		Message:  p.Message,
		IsEdited: p.IsEdited,
		Forum:    p.Forum,
		Thread:   p.Thread,
		Created:  p.Created,
	}
}

func DbToJsonModel(p *Post) *models.Post {
	return &models.Post{
		Id:       p.Id,
		Parent:   p.Parent,
		Author:   p.Author,
		Message:  p.Message,
		IsEdited: p.IsEdited,
		Forum:    p.Forum,
		Thread:   p.Thread,
		Created:  p.Created,
	}
}
