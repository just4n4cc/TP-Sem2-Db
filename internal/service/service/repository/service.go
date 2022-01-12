package repository

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type Status struct {
	Id     int   `db:"id"`
	User   int32 `db:"users"`
	Forum  int32 `db:"forums"`
	Thread int32 `db:"threads"`
	Post   int64 `db:"posts"`
}

func JsonToDbModel(p *models.Status) *Status {
	return &Status{
		User:   p.User,
		Forum:  p.Forum,
		Thread: p.Thread,
		Post:   p.Post,
	}
}
func DbToJsonModel(p *Status) *models.Status {
	return &models.Status{
		User:   p.User,
		Forum:  p.Forum,
		Thread: p.Thread,
		Post:   p.Post,
	}
}
