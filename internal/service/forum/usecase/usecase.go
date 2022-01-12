package usecase

import (
	"github.com/gofrs/uuid"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/service/forum"
	thread "github.com/just4n4cc/tp-sem2-db/internal/service/thread"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
)

type Usecase struct {
	repositoryForum  forum.Repository
	repositoryThread thread.Repository
}

func NewUseCase(repositoryForum forum.Repository, repositoryThread thread.Repository) *Usecase {
	return &Usecase{
		repositoryForum:  repositoryForum,
		repositoryThread: repositoryThread,
	}
}

func (a *Usecase) ForumCreate(forum *models.Forum) (*models.Forum, error) {
	if forum.Posts != 0 || forum.Threads != 0 || forum.Slug == "" || forum.User == "" || forum.Title == "" {
		return nil, models.ModelFieldError
	}
	return a.repositoryForum.ForumCreate(forum)
}
func (a *Usecase) ForumGet(slug string) (*models.Forum, error) {
	if slug == "" {
		return nil, models.ModelFieldError
	}
	return a.repositoryForum.ForumGet(slug)
}
func (a *Usecase) ForumThreadCreate(thread *models.Thread) (*models.Thread, error) {

	if thread.Title == "" || thread.Forum == "" || thread.Author == "" || thread.Message == "" {
		log.Debug("here?")
		return nil, models.ModelFieldError
	}
	u, err := uuid.NewV4()
	log.Debug("here?")
	if err != nil {
		return nil, err
	}
	log.Debug("mb?")
	thread.Slug = u.String()
	log.Debug("tyt?")
	// TODO initialize created
	return a.repositoryThread.ThreadCreate(thread)
}
func (a *Usecase) ForumUsers(slug string, so *models.SortOptions) ([]*models.User, error) {
	if slug == "" {
		return nil, models.ModelFieldError
	}
	_, err := a.repositoryForum.ForumGet(slug)
	if err != nil {
		return nil, err
	}
	users, err := a.repositoryForum.ForumUsers(slug, so)
	// TODO some sort
	return users, err
}
func (a *Usecase) ForumThreads(slug string, so *models.SortOptions) ([]*models.Thread, error) {
	if slug == "" {
		return nil, models.ModelFieldError
	}
	_, err := a.repositoryForum.ForumGet(slug)
	if err != nil {
		return nil, err
	}
	threads, err := a.repositoryThread.ThreadsByForum(slug, so)
	// TODO some sort
	return threads, err
}
