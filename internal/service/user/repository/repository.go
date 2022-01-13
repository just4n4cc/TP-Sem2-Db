package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
	"strconv"
	"strings"
)

const (
	logMessage            = "service:user:repository:"
	selectAlreadyExisting = `select * from tpdb."User"
		where nickname = $1 or email = $2`
	userCreate = `insert into tpdb."User"
		(nickname, fullname, about, email)
		values($1, $2, $3, $4)`
	userProfileGet = `select * from tpdb."User"
		where nickname = $1`
	userProfileUpdateBegin = `update tpdb."User" set `
	userProfileUpdateEnd   = `where nickname = $1
		returning *`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) UserCreate(u *models.User) ([]*models.User, error) {
	message := logMessage + "UserCreate:"
	log.Debug(message + "started")
	user := JsonToDbModel(u)
	query := userCreate
	rows, err := s.db.Queryx(query, user.Nickname, user.Fullname, user.About, user.Email)
	if rows != nil {
		e := rows.Close()
		if e != nil {
			log.Error(message, e)
		}
	}
	if err == nil {
		log.Success(message)
		return nil, nil
	}
	err = utils.TranslateDbError(err)
	if err != models.AlreadyExistsError {
		log.Error(message, err)
		return nil, err
	}

	var users []User
	err = s.db.Select(&users, selectAlreadyExisting, user.Nickname, user.Email)
	if err != nil {
		log.Error(message, err)
		return nil, err
	}
	//if len(users) == 0 {
	//	err = models.UnexpectedDbBehavior
	//	log.Error(message, err)
	//	return nil, err
	//}
	//var us []*models.User
	//for _, u := range users {
	//	us = append(us, DbToJsonModel(&u))
	//}
	log.Success(message)
	//return us, models.AlreadyExistsError
	return DbArrayToJsonModel(users), models.AlreadyExistsError
}

func (s *Repository) UserProfileGet(nickname string) (*models.User, error) {
	message := logMessage + "UserProfileGet:"
	log.Debug(message + "started")
	query := userProfileGet
	user := new(User)
	err := s.db.Get(user, query, nickname)
	//log.Error(message, err)
	if err == nil {
		log.Success(message)
		return DbToJsonModel(user), nil
	}
	err = utils.TranslateDbError(err)
	return nil, err
}

func (s *Repository) UserProfileUpdate(u *models.User) (*models.User, error) {
	message := logMessage + "UserProfileUpdate:"
	log.Debug(message + "started")
	user := JsonToDbModel(u)
	query := userProfileUpdateBegin
	//set fullname = $2, about = $3, email = $4
	num := 2
	var args []interface{}
	args = append(args, user.Nickname)
	if user.Fullname != "" {
		query += "fullname = $" + strconv.Itoa(num) + ", "
		args = append(args, user.Fullname)
		num++
	}
	if user.About != "" {
		query += "about = $" + strconv.Itoa(num) + ", "
		args = append(args, user.About)
		num++
	}
	if user.Email != "" {
		query += "email = $" + strconv.Itoa(num) + ", "
		args = append(args, user.Email)
		num++
	}
	query = strings.TrimSuffix(query, ", ")
	query += " " + userProfileUpdateEnd
	log.Debug(message + "query = " + query)
	err := s.db.Get(user, query, args...)
	if err == nil {
		log.Success(message)
		//return nil, nil
		return DbToJsonModel(user), nil
	}
	err = utils.TranslateDbError(err)
	return nil, err
}
