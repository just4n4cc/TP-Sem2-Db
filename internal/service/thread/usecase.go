package thread

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type UseCase interface {
	ThreadCreatePosts(slugOrId string, posts []*models.Post) ([]*models.Post, error)
	ThreadBySlugOrId(slugOrId string) (*models.Thread, error)
	ThreadUpdate(slugOrId string, thread *models.Thread) (*models.Thread, error)
	ThreadPosts(slugOrId string, so *models.SortOptions) ([]*models.Post, error)
	ThreadVote(slugOrId string, vote *models.Vote) (*models.Thread, error)
}
