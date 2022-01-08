package usecase

import (
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/service/forum"
	repoThread "github.com/just4n4cc/tp-sem2-db/internal/service/thread/repository"
)

type Usecase struct {
	repositoryForum  forum.Repository
	reporitoryThread repoThread.Repository
}

func NewUseCase(repositoryForum forum.Repository, repositoryThread repoThread.Repository) *Usecase {
	return &Usecase{
		repositoryForum: repositoryForum,
	}
}

func (a *Usecase) ForumCreate(forum *models.Forum) error {
	if forum.Posts != 0 || forum.Threads != 0 || forum.Slug == "" || forum.User == "" || forum.Title == "" {
		return models.ModelFieldError
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
		return nil, models.ModelFieldError
	}
	return a.repositoryForum.ForumThreadCreate(thread)
}
func (a *Usecase) ForumUsers(slug string, so *models.SortOptions) ([]*models.User, error) {
	if slug == "" {
		return nil, models.ModelFieldError
	}
	users, err := a.repositoryForum.ForumUsers(slug, so)
	// TODO some sort
	return users, err
}
func (a *Usecase) ForumThreads(slug string, so *models.SortOptions) ([]*models.Thread, error) {
	if slug == "" {
		return nil, models.ModelFieldError
	}
	threads, err := a.reporitoryThread.ThreadsByForum(slug, so)
	// TODO some sort
	return threads, err
}
