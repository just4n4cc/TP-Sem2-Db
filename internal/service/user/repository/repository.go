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
	selectAlreadyExisting = `select * from "User"
		where nickname = $1 or email = $2`
	userCreate = `insert into "User"
		(nickname, fullname, about, email)
		values($1, $2, $3, $4)`
	userProfileGet = `select * from "User"
		where nickname = $1`
	userProfileUpdateBegin = `update "User" set `
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
	log.Success(message)
	return DbArrayToJsonModel(users), models.AlreadyExistsError
}

func (s *Repository) UserProfileGet(nickname string) (*models.User, error) {
	message := logMessage + "UserProfileGet:"
	log.Debug(message + "started")
	query := userProfileGet
	user := new(User)
	err := s.db.Get(user, query, nickname)
	if err == nil {
		log.Success(message)
		return DbToJsonModel(user), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}

func (s *Repository) UserProfileUpdate(u *models.User) (*models.User, error) {
	message := logMessage + "UserProfileUpdate:"
	log.Debug(message + "started")
	user := JsonToDbModel(u)
	query := userProfileUpdateBegin
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
	err := s.db.Get(user, query, args...)
	if err == nil {
		log.Success(message)
		return DbToJsonModel(user), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
