package usecase

import (
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/service/forum"
	"github.com/just4n4cc/tp-sem2-db/internal/service/post"
	"github.com/just4n4cc/tp-sem2-db/internal/service/thread"
	"github.com/just4n4cc/tp-sem2-db/internal/service/user"
	"strconv"
	"strings"
)

type UseCase struct {
	repositoryPost   post.Repository
	repositoryForum  forum.Repository
	repositoryThread thread.Repository
	repositoryUser   user.Repository
}

func NewUseCase(repositoryPost post.Repository, repositoryForum forum.Repository, repositoryThread thread.Repository, repositoryUser user.Repository) *UseCase {
	return &UseCase{
		repositoryPost:   repositoryPost,
		repositoryForum:  repositoryForum,
		repositoryThread: repositoryThread,
		repositoryUser:   repositoryUser,
	}
}

func (a *UseCase) PostGet(id string, related []string) (*models.Post, *models.Forum, *models.Thread, *models.User, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, nil, nil, nil, models.ModelFieldError
	}
	p, err := a.repositoryPost.PostGet(i)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	relatedModels := ""
	for _, r := range related {
		relatedModels += r
	}
	var f *models.Forum
	var t *models.Thread
	var u *models.User
	f = nil
	t = nil
	u = nil
	if strings.Contains(relatedModels, "user") {
		u, err = a.repositoryUser.UserProfileGet(p.Author)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if strings.Contains(relatedModels, "thread") {
		t, err = a.repositoryThread.ThreadById(p.Thread)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if strings.Contains(relatedModels, "forum") {
		f, err = a.repositoryForum.ForumGet(p.Forum)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	return p, f, t, u, nil
}
func (a *UseCase) PostUpdate(post *models.Post) (*models.Post, error) {
	if post.Message == "" {
		return a.repositoryPost.PostGet(post.Id)
	}
	return a.repositoryPost.PostUpdate(post)
}
