package response

import (
	"encoding/json"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"net/http"
)

//func GetStatus(err error, message string) int {
func GetStatus(err error) int {
	if err == models.NotFoundError {
		return http.StatusNotFound
	}
	if err == models.AlreadyExistsError {
		return http.StatusConflict
	}
	return http.StatusBadGateway
}

func UnknownError(w *http.ResponseWriter) {
	////logger.Error(message, err)
	(*w).WriteHeader(http.StatusInternalServerError)
	return
}

//func SetBody(w *http.ResponseWriter, object interface{}, message string) {
func SetBody(w *http.ResponseWriter, object interface{}) {
	if object == nil {
		////logger.Success(message)
		return
	}
	body, err := json.Marshal(object)
	if err != nil {
		UnknownError(w)
		//UnknownError(w, err, message)
		////logger.Error(message, err)
		return
	}
	b := string(body)
	if b == "null" {
		switch object.(type) {
		case []*models.Thread, []*models.Forum, []*models.User, []*models.Post:
			b = "[]"
		default:
			b = "{}"
		}
		body = []byte(b)
	}
	_, err = (*w).Write(body)
	if err != nil {
		UnknownError(w)
		////logger.Error(message, err)
		return
	}
	////logger.Success(message)
}
