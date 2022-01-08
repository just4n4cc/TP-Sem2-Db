package utils

import (
	"errors"
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/pkg/logger"
	"net/http"
	"strconv"
	"strings"
)

const logMessage = "utils:"

func InitDb() (*sql.DB, error) {
	message := logMessage + "InitDb:"

	db, err := sql.Connect("postgres",
		"host=127.0.0.1 port=5000 user=postgres dbname=postgres password=password sslmode=disable")
	//"user=postgres dbname=tpdb sslmode=disable")
	if err != nil {
		logger.Error(message+"err = ", err)
		return nil, err
	}
	logger.Info(message + "[SUCCESS]")
	return db, nil
}

func GetSortOptionsFromRequest(r *http.Request) (*models.SortOptions, error) {
	so := new(models.SortOptions)
	q := r.URL.Query()
	const limit = "limit"
	const since = "since"
	const sort = "sort"
	const desc = "desc"

	if len(q[limit]) > 0 {
		p, err := strconv.ParseInt(q[limit][0], 10, 32)
		if err != nil {
			return nil, err
		}
		so.Limit = int32(p)
	} else {
		so.Limit = 100
	}

	if len(q[since]) > 0 {
		p, err := strconv.ParseInt(q[since][0], 10, 64)
		if err != nil {
			return nil, err
		}
		so.Since = p
	}

	if len(q[sort]) > 0 {
		p := q[sort][0]
		if p == "flat" || p == "tree" || p == "parent_tree" {
			so.Sort = p
		} else {
			return nil, errors.New("unexpected sort type")
		}
	} else {
		so.Sort = "flat"
	}

	if len(q[desc]) > 0 {
		p, err := strconv.ParseBool(q[since][0])
		if err != nil {
			return nil, err
		}
		so.Desc = p
	}
	return so, nil
}

func TranslateDbError(err error) error {
	if err == nil {
		return nil
	}
	if strings.HasPrefix(err.Error(), "pq: duplicate key value violates unique constraint") {
		return models.AlreadyExistsError
	}
	if err.Error() == "sql: no rows in result set" {
		return models.NotFoundError
	}
	return err
}
