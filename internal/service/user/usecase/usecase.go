package usecase

import (
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/service/user"
)

type Usecase struct {
	repository user.Repository
}

func NewUseCase(repository user.Repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (a *Usecase) UserCreate(u *models.User) ([]*models.User, error) {
	if u.Nickname == "" || u.Fullname == "" || u.Email == "" {
		return nil, models.ModelFieldError
	}

	return a.repository.UserCreate(u)
}

func (a *Usecase) UserProfileGet(nickname string) (*models.User, error) {
	if nickname == "" {
		return nil, models.ModelFieldError
	}
	return a.repository.UserProfileGet(nickname)
}

func (a *Usecase) UserProfileUpdate(u *models.User) (*models.User, error) {
	if u.About == "" && u.Fullname == "" && u.Email == "" {
		return a.UserProfileGet(u.Nickname)
	}

	if u.Nickname == "" {
		return nil, models.ModelFieldError
	}
	return a.repository.UserProfileUpdate(u)
}
