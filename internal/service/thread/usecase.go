package thread

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type UseCase interface {
	ThreadCreate(slugOrId string, posts []*models.Post) ([]*models.Post, error)
	ThreadGet(slugOrId string) (*models.Thread, error)
	ThreadUpdate(thread *models.Thread) (*models.Thread, error)
	ThreadPosts(slugOrId string, so *models.SortOptions) ([]*models.Post, error)
	ThreadVote(slugOrId string, vote *models.Vote) (*models.Thread, error)
}
