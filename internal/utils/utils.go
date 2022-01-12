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
		"host=127.0.0.1 port=5432 user=just4n4cc dbname=postgres password=password sslmode=disable")
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
		so.Since = q[since][0]
		//p, err := strconv.ParseInt(q[since][0], 10, 64)
		//if err != nil {
		//	return nil, err
		//}
		//so.Since = p
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
		p, err := strconv.ParseBool(q[desc][0])
		if err != nil {
			return nil, err
		}
		so.Desc = p
	}
	return so, nil
}

func SortOptionsToSubquery(so *models.SortOptions, field string) string {
	if so == nil {
		return ""
	}
	res := ""
	sign := ""
	order := ""
	if !so.Desc {
		sign = "<="
		order = "desc"
	} else {
		sign = ">="
	}
	if so.Since != "" {
		res += "and " + field + "" + sign
	}
	res += " order by " + field + " " + order + " limit " + strconv.Itoa(int(so.Limit))
	return res
}

func TranslateDbError(err error) error {
	if err == nil {
		return nil
	}
	s := err.Error()
	if strings.HasPrefix(s, "pq: duplicate key value violates unique constraint") {
		return models.AlreadyExistsError
	}
	if s == "sql: no rows in result set" || strings.Contains(s, "violates foreign key constraint") {
		return models.NotFoundError
	}
	return err
}

type Set struct {
	s []interface{}
}

func NewSet(args []interface{}) *Set {
	set := new(Set)
	for _, arg := range args {
		if !set.Contains(arg) {
			set.s = append(set.s, arg)
		}
	}
	return set
}

func (s *Set) Contains(el interface{}) bool {
	for _, e := range s.s {
		if e == el {
			return true
		}
	}
	return false
}

func (s *Set) Add(el interface{}) {
	if !s.Contains(el) {
		s.s = append(s.s, el)
	}
}

func (s *Set) Length() int {
	return len(s.s)
}
