package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/response"
	"github.com/just4n4cc/tp-sem2-db/internal/service/forum"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	"github.com/just4n4cc/tp-sem2-db/pkg/logger"
	"net/http"
)

const logMessage = "service:forum:delivery:"
const slug = "slug"

type Delivery struct {
	useCase forum.UseCase
}

func NewDelivery(useCase forum.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func (h *Delivery) ForumCreate(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ForumCreate:"
	logger.Debug(message + "started")

	var f = new(models.Forum)
	err := json.NewDecoder(r.Body).Decode(f)
	if err != nil {
		response.UnknownError(&w, err, message)
		return
	}

	err = h.useCase.ForumCreate(f)
	if err != nil {
		if err == models.AlreadyExistsError {
			w.WriteHeader(response.GetStatus(err, message))
			response.SetBody(&w, f, message)
			return
		} else if err == models.NotFoundError {
			w.WriteHeader(response.GetStatus(err, message))
			e := models.Error{Message: err.Error()}
			response.SetBody(&w, e, message)
			return
		} else {
			response.UnknownError(&w, err, message)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	response.SetBody(&w, f, message)
}

func (h *Delivery) ForumGet(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ForumGet:"
	logger.Debug(message + "started")
	vars := mux.Vars(r)
	logger.Debug(message + "slug = " + vars[slug])

	f, err := h.useCase.ForumGet(vars[slug])
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
	response.SetBody(&w, f, message)
}

func (h *Delivery) ForumThreadCreate(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ForumThreadCreate:"
	logger.Debug(message + "started")
	vars := mux.Vars(r)
	logger.Debug(message + "slug = " + vars[slug])

	var t = new(models.Thread)
	err := json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		response.UnknownError(&w, err, message)
		return
	}
	t.Forum = vars[slug]

	t, err = h.useCase.ForumThreadCreate(t)

	if err != nil {
		if err == models.AlreadyExistsError {
			w.WriteHeader(response.GetStatus(err, message))
			response.SetBody(&w, t, message)
			return
		} else if err == models.NotFoundError {
			w.WriteHeader(response.GetStatus(err, message))
			e := models.Error{Message: err.Error()}
			response.SetBody(&w, e, message)
			return
		} else {
			response.UnknownError(&w, err, message)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	response.SetBody(&w, t, message)
}

func (h *Delivery) ForumUsers(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ForumUsers:"
	logger.Debug(message + "started")
	vars := mux.Vars(r)
	logger.Debug(message + "slug = " + vars[slug])
	so, err := utils.GetSortOptionsFromRequest(r)
	if err != nil {
		response.UnknownError(&w, err, message)
	}
	logger.Debug(message+"sort options = ", so)

	users, err := h.useCase.ForumUsers(vars[slug], so)
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
	response.SetBody(&w, users, message)
}

func (h *Delivery) ForumThreads(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ForumThreads:"
	logger.Debug(message + "started")
	vars := mux.Vars(r)
	logger.Debug(message + "slug = " + vars[slug])
	so, err := utils.GetSortOptionsFromRequest(r)
	if err != nil {
		response.UnknownError(&w, err, message)
	}
	logger.Debug(message+"sort options = ", so)

	threads, err := h.useCase.ForumUsers(vars[slug], so)
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
	response.SetBody(&w, threads, message)
}
