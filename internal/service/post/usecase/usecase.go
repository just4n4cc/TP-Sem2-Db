package usecase

import (
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/service/forum"
	"github.com/just4n4cc/tp-sem2-db/internal/service/post"
	"github.com/just4n4cc/tp-sem2-db/internal/service/thread"
	"github.com/just4n4cc/tp-sem2-db/internal/service/user"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
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
	post, err := a.repositoryPost.PostGet(i)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	relatedModels := ""
	log.Debug("here")
	for _, r := range related {
		relatedModels += r
	}
	log.Debug("here")
	var f *models.Forum
	var t *models.Thread
	var u *models.User
	f = nil
	t = nil
	u = nil
	if strings.Contains(relatedModels, "user") {
		u, err = a.repositoryUser.UserProfileGet(post.Author)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if strings.Contains(relatedModels, "thread") {
		t, err = a.repositoryThread.ThreadById(post.Thread)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if strings.Contains(relatedModels, "forum") {
		f, err = a.repositoryForum.ForumGet(post.Forum)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	return post, f, t, u, nil
}
func (a *UseCase) PostUpdate(post *models.Post) (*models.Post, error) {
	if post.Message == "" {
		return nil, models.ModelFieldError
	}
	return a.repositoryPost.PostUpdate(post)
}

//func (a *UseCase) PostCreate(p *models.Post) (*models.Post, error) {
//	if p.Message == "" || p.Forum == "" || p.Author
//}
//
//func (a *UseCase) PostGet(id string, related []string) (*models.Post, *models.User, *models.Thread, *models.Forum, error) {
//
//}
//func (a *UseCase) PostUpdate(*models.Post) (*models.Post, error) {
//
//}
