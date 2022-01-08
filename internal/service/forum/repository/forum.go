package repository

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type Forum struct {
	Id      int    `db:"id"`
	Title   string `db:"title"`
	User    string `db:"user"`
	Slug    string `db:"slug"`
	Posts   int64  `db:"posts"`
	Threads int32  `db:"threads"`
}

func jsonToDbModel(f *models.Forum) *Forum {
	return &Forum{
		Id:      0,
		Title:   f.Title,
		User:    f.User,
		Slug:    f.Slug,
		Posts:   f.Posts,
		Threads: f.Threads,
	}
}

func dbToJsonModel(f *Forum) *models.Forum {
	return &models.Forum{
		Title:   f.Title,
		User:    f.User,
		Slug:    f.Slug,
		Posts:   f.Posts,
		Threads: f.Threads,
	}
}
