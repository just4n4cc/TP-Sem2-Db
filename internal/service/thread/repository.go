package thread

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type Repository interface {
	ThreadsByForum(forum string, so *models.SortOptions) ([]*models.Thread, error)
	ThreadCreate(t *models.Thread, slug string) (*models.Thread, error)
	ThreadBySlug(slugOrId string) (*models.Thread, error)
	ThreadById(id int32) (*models.Thread, error)
	ThreadUpdateById(t *models.Thread) (*models.Thread, error)
	ThreadUpdateBySlug(t *models.Thread) (*models.Thread, error)
	ThreadVote(t *models.Thread, v *models.Vote) (*models.Thread, error)
	//ThreadVoteBySlug(slug string, vote *models.Vote) (*models.Thread, error)
}
