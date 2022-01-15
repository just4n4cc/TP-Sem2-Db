package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/response"
	"github.com/just4n4cc/tp-sem2-db/internal/service/thread"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
	"net/http"
)

const logMessage = "service:thread:delivery:"
const slugOrId = "slug_or_id"

type Delivery struct {
	useCase thread.UseCase
}

func NewDelivery(useCase thread.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func (h *Delivery) ThreadCreate(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ThreadCreate:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	log.Debug(message + "slug or id = " + vars[slugOrId])

	var posts []*models.Post
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		response.UnknownError(&w)
		return
	}

	posts, err = h.useCase.ThreadCreatePosts(vars[slugOrId], posts)
	if err != nil {
		if err != models.AlreadyExistsError && err != models.NotFoundError {
			response.UnknownError(&w)
			return
		}

		w.WriteHeader(response.GetStatus(err))
		e := models.Error{Message: err.Error()}
		response.SetBody(&w, e)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response.SetBody(&w, posts)
}

func (h *Delivery) ThreadGet(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ThreadGet:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	log.Debug(message + "slug or id = " + vars[slugOrId])

	t, err := h.useCase.ThreadBySlugOrId(vars[slugOrId])
	if err != nil {
		if err != models.NotFoundError {
			response.UnknownError(&w)
			return
		}

		w.WriteHeader(response.GetStatus(err))
		e := models.Error{Message: err.Error()}
		response.SetBody(&w, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.SetBody(&w, t)
}

func (h *Delivery) ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ThreadUpdate:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	log.Debug(message + "slug or id = " + vars[slugOrId])

	var t = new(models.Thread)
	err := json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		response.UnknownError(&w)
		return
	}

	t, err = h.useCase.ThreadUpdate(vars[slugOrId], t)
	if err != nil {
		if err != models.NotFoundError {
			response.UnknownError(&w)
			return
		}

		w.WriteHeader(response.GetStatus(err))
		e := models.Error{Message: err.Error()}
		response.SetBody(&w, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.SetBody(&w, t)
}

func (h *Delivery) ThreadPosts(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ThreadPosts:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	log.Debug(message + "slug or id = " + vars[slugOrId])
	so, err := utils.GetSortOptionsFromRequest(r)
	if err != nil {
		response.UnknownError(&w)
	}
	log.Debug(message+"sort options = ", so)

	posts, err := h.useCase.ThreadPosts(vars[slugOrId], so)
	if err != nil {
		if err != models.NotFoundError {
			response.UnknownError(&w)
			return
		}

		w.WriteHeader(response.GetStatus(err))
		e := models.Error{Message: err.Error()}
		response.SetBody(&w, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.SetBody(&w, posts)
}

func (h *Delivery) ThreadVote(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "ThreadVote:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	log.Debug(message + "slug or id = " + vars[slugOrId])

	var v = new(models.Vote)
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		response.UnknownError(&w)
		return
	}

	p, err := h.useCase.ThreadVote(vars[slugOrId], v)
	if err != nil {
		if err != models.NotFoundError {
			response.UnknownError(&w)
			return
		}

		w.WriteHeader(response.GetStatus(err))
		e := models.Error{Message: err.Error()}
		response.SetBody(&w, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.SetBody(&w, p)
}
