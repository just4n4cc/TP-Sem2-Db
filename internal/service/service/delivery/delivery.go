package delivery

import (
	"github.com/just4n4cc/tp-sem2-db/internal/response"
	"github.com/just4n4cc/tp-sem2-db/internal/service/service"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
	"net/http"
)

const logMessage = "service:service:delivery:"

type Delivery struct {
	useCase service.UseCase
}

func NewDelivery(useCase service.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func (h *Delivery) ServiceClear(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ServiceClear:"
	log.Debug(message + "started")

	err := h.useCase.ServiceClear()
	if err != nil {
		response.UnknownError(&w)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.SetBody(&w, nil)
}

func (h *Delivery) ServiceStatus(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ServiceStatus:"
	log.Debug(message + "started")

	s, err := h.useCase.ServiceStatus()
	if err != nil {
		response.UnknownError(&w)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.SetBody(&w, s)
}
