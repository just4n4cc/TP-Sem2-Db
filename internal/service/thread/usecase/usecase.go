package usecase

import (
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/service/post"
	"github.com/just4n4cc/tp-sem2-db/internal/service/thread"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	"strconv"
	"time"
)

type Usecase struct {
	repositoryThread thread.Repository
	repositoryPost   post.Repository
}

func NewUseCase(repositoryThread thread.Repository, repositoryPost post.Repository) *Usecase {
	return &Usecase{
		repositoryThread: repositoryThread,
		repositoryPost:   repositoryPost,
	}
}
func (a *Usecase) ThreadCreatePosts(slugOrId string, posts []*models.Post) ([]*models.Post, error) {
	ps, err := a.ThreadPosts(slugOrId, nil)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return posts, nil
	}
	idSet := new(utils.Set)
	idSet.Add(int64(0))

	var forum string
	var threadId int32
	if len(ps) == 0 {
		t, err := a.ThreadBySlugOrId(slugOrId)
		if err != nil {
			return nil, err
		}
		forum = t.Forum
		threadId = t.Id
	} else {
		forum = ps[0].Forum
		threadId = ps[0].Thread
	}

	for _, p := range ps {
		idSet.Add(p.Id)
	}

	if err != nil {
		return nil, err
	}

	oclock := time.Now()
	for _, p := range posts {
		p.Forum = forum
		p.Thread = threadId
		p.Created = oclock
		if p.Author == "" || p.Message == "" {
			return nil, models.ModelFieldError
		}
		if !idSet.Contains(p.Parent) {
			return nil, models.AlreadyExistsError
		}
	}
	return a.repositoryPost.PostsCreate(posts)
}
func (a *Usecase) ThreadBySlugOrId(slugOrId string) (*models.Thread, error) {
	if slugOrId == "" {
		return nil, models.ModelFieldError
	}
	id, err := strconv.ParseInt(slugOrId, 10, 32)
	if err == nil {
		return a.repositoryThread.ThreadById(int32(id))
	}
	return a.repositoryThread.ThreadBySlug(slugOrId)
}

func (a *Usecase) ThreadUpdate(slugOrId string, thread *models.Thread) (*models.Thread, error) {
	if thread.Message == "" && thread.Title == "" {
		return a.ThreadBySlugOrId(slugOrId)
	}
	if slugOrId == "" {
		return nil, models.ModelFieldError
	}
	//if slugOrId == "" || thread.Message == "" || thread.Title == "" {
	//	return nil, models.ModelFieldError
	//}
	id, err := strconv.ParseInt(slugOrId, 10, 32)
	if err == nil {
		thread.Id = int32(id)
		return a.repositoryThread.ThreadUpdateById(thread)
	}
	thread.Slug = slugOrId
	return a.repositoryThread.ThreadUpdateBySlug(thread)
}

func (a *Usecase) ThreadPosts(slugOrId string, so *models.SortOptions) ([]*models.Post, error) {
	t, err := a.ThreadBySlugOrId(slugOrId)
	if err != nil {
		return nil, err
	}
	return a.repositoryPost.PostsByThread(t.Id, so)
}
func (a *Usecase) ThreadVote(slugOrId string, vote *models.Vote) (*models.Thread, error) {
	if slugOrId == "" {
		return nil, models.ModelFieldError
	}
	if vote.Voice != -1 && vote.Voice != 1 {
		return nil, models.ModelFieldError
	}

	t, err := a.ThreadBySlugOrId(slugOrId)
	if err != nil {
		return nil, err
	}
	return a.repositoryThread.ThreadVote(t, vote)
}
