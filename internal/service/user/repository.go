package user

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type Repository interface {
	UserCreate(u *models.User) ([]*models.User, error)
	UserProfileGet(nickname string) (*models.User, error)
	UserProfileUpdate(u *models.User) (*models.User, error)
}
