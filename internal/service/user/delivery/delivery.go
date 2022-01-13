package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/response"
	"github.com/just4n4cc/tp-sem2-db/internal/service/user"
	"github.com/just4n4cc/tp-sem2-db/pkg/logger"
	"net/http"
)

const logMessage = "service:user:delivery:http:"
const nickname = "nickname"

type Delivery struct {
	useCase user.UseCase
}

func NewDelivery(useCase user.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func (h *Delivery) UserCreate(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UserCreate:"
	logger.Debug(message + "started")
	vars := mux.Vars(r)
	logger.Debug(message + "nickname = " + vars[nickname])

	var u = new(models.User)
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		response.UnknownError(&w, err, message)
		return
	}
	u.Nickname = vars[nickname]

	users, err := h.useCase.UserCreate(u)
	if err != nil {
		if err != models.AlreadyExistsError {
			response.UnknownError(&w, err, message)
			return
		}

		w.WriteHeader(response.GetStatus(err, message))
		response.SetBody(&w, users, message)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response.SetBody(&w, u, message)
}

func (h *Delivery) UserProfileGet(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UserProfileGet:"
	logger.Debug(message + "started")
	vars := mux.Vars(r)
	logger.Debug(message + "nickname = " + vars[nickname])

	u, err := h.useCase.UserProfileGet(vars[nickname])
	if err != nil {
		if err != models.NotFoundError {
			response.UnknownError(&w, err, message)
			return
		}

		w.WriteHeader(response.GetStatus(err, message))
		e := models.Error{Message: err.Error()}
		response.SetBody(&w, e, message)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.SetBody(&w, u, message)
}

func (h *Delivery) UserProfileUpdate(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UserProfileUpdate:"
	logger.Debug(message + "started")
	vars := mux.Vars(r)
	logger.Debug(message + "nickname = " + vars[nickname])

	var u = new(models.User)
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		response.UnknownError(&w, err, message)
		return
	}
	u.Nickname = vars[nickname]

	u, err = h.useCase.UserProfileUpdate(u)
	//err = h.useCase.UserProfileUpdate(u)
	if err != nil {
		if err != models.AlreadyExistsError && err != models.NotFoundError {
			response.UnknownError(&w, err, message)
			return
		}

		w.WriteHeader(response.GetStatus(err, message))
		e := models.Error{Message: err.Error()}
		response.SetBody(&w, e, message)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.SetBody(&w, u, message)
}
