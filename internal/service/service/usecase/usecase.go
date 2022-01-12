package usecase

import (
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/service/service"
)

type UseCase struct {
	repository service.Repository
}

func NewUseCase(repositoryService service.Repository) *UseCase {
	return &UseCase{
		repository: repositoryService,
	}
}

func (a *UseCase) ServiceClear() error {
	return a.repository.ServiceClear()
}

func (a *UseCase) ServiceStatus() (*models.Status, error) {
	return a.repository.ServiceStatus()
}
