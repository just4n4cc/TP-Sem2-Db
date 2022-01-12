package forum

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type Repository interface {
	ForumCreate(f *models.Forum) (*models.Forum, error)
	ForumGet(slug string) (*models.Forum, error)
	ForumUsers(slug string, so *models.SortOptions) ([]*models.User, error)
}
