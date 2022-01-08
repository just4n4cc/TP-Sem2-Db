package post

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type UseCase interface {
	PostGet(id string, related []string) (*models.Post, *models.User, *models.Thread, *models.Forum, error)
	PostUpdate(*models.Post) (*models.Post, error)
}
