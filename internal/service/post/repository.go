package post

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type Repository interface {
	PostGet(id int64) (*models.Post, error)
	PostUpdate(p *models.Post) (*models.Post, error)
	PostsCreate(ps []*models.Post) ([]*models.Post, error)
	PostsByThread(id int32, so *models.SortOptions) ([]*models.Post, error)
}
