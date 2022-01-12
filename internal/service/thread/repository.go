package thread

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type Repository interface {
	ThreadsByForum(forum string, so *models.SortOptions) ([]*models.Thread, error)
	ThreadCreate(t *models.Thread) (*models.Thread, error)
	ThreadBySlug(slugOrId string) (*models.Thread, error)
	ThreadById(id int32) (*models.Thread, error)
	ThreadUpdateById(t *models.Thread) (*models.Thread, error)
	ThreadUpdateBySlug(t *models.Thread) (*models.Thread, error)
	ThreadVoteById(id int32, vote *models.Vote) (*models.Thread, error)
	ThreadVoteBySlug(slug string, vote *models.Vote) (*models.Thread, error)
}
