package service

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type UseCase interface {
	ServiceClear() error
	ServiceStatus() (*models.Status, error)
}
