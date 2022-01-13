package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/response"
	"github.com/just4n4cc/tp-sem2-db/internal/service/post"
	"github.com/just4n4cc/tp-sem2-db/pkg/logger"
	"net/http"
	"strconv"
)

const logMessage = "service:post:delivery:http:"
const id = "id"
const related = "related"

type Delivery struct {
	useCase post.UseCase
}

func NewDelivery(useCase post.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

type PostGetResponse struct {
	Post   *models.Post   `json:"post,omitempty"`
	Author *models.User   `json:"author,omitempty"`
	Thread *models.Thread `json:"thread,omitempty"`
	Forum  *models.Forum  `json:"forum,omitempty"`
}

func (h *Delivery) PostGet(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "PostGet:"
	logger.Debug(message + "started")
	vars := mux.Vars(r)
	logger.Debug(message + "id = " + vars[id])
	query := r.URL.Query()
	logger.Debug(message+"related = ", query[related])

	p, f, t, u, err := h.useCase.PostGet(vars[id], query[related])
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
	b := &PostGetResponse{
		Post: p,
	}
	if f != nil {
		b.Forum = f
	}
	if t != nil {
		b.Thread = t
	}
	if u != nil {
		b.Author = u
	}
	response.SetBody(&w, b, message)
}

func (h *Delivery) PostUpdate(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "PostUpdate:"
	logger.Debug(message + "started")
	vars := mux.Vars(r)
	logger.Debug(message + "id = " + vars[id])

	var p = new(models.Post)
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		response.UnknownError(&w, err, message)
		return
	}
	p.Id, err = strconv.ParseInt(vars[id], 10, 64)
	if err != nil {
		response.UnknownError(&w, err, message)
		return
	}

	p, err = h.useCase.PostUpdate(p)
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
	response.SetBody(&w, p, message)
}
