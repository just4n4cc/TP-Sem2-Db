package forum

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type UseCase interface {
	ForumCreate(forum *models.Forum) (*models.Forum, error)
	ForumGet(slug string) (*models.Forum, error)
	ForumThreadCreate(thread *models.Thread) (*models.Thread, error)
	ForumUsers(slug string, so *models.SortOptions) ([]*models.User, error)
	ForumThreads(slug string, so *models.SortOptions) ([]*models.Thread, error)
}
