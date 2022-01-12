package post

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type UseCase interface {
	PostGet(id string, related []string) (*models.Post, *models.Forum, *models.Thread, *models.User, error)
	PostUpdate(post *models.Post) (*models.Post, error)
}
