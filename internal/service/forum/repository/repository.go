package repository

import (
	"errors"
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
)

const (
	logMessage = "service:forum:repository:"
	//selectAlreadyExist = `select * from tpdb."Forum"
	//	where slug = $1`
	forumCreate = `insert into tpdb."Forum"
		(title, user, slug)
		values($1, $2, $3)`
	userProfileGet = `select * from tpdb."User"
		where nickname = $1`
	userProfileUpdate = `update tpdb."User"
		set fullname = $2, about = $3, email = $4
		where nickname = $1
		returning id`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) ForumGet(slug string) (*models.Forum, error) {
	message := logMessage + "ForumGet:"
	log.Debug(message + "started")
	forum := new(Forum)
	//err := s.db.Select(forum, selectAlreadyExist, slug)
	err := errors.New("lala")
	if err == nil {
		log.Success(message)
		return dbToJsonModel(forum), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}

func (s *Repository) ForumCreate(f *models.Forum) (*models.Forum, error) {
	message := logMessage + "ForumCreate:"
	log.Debug(message + "started")
	forum := jsonToDbModel(f)
	query := forumCreate
	_, err := s.db.Queryx(query, forum.Title, forum.User, forum.Slug)
	if err == nil {
		log.Success(message)
		return f, nil
	}
	err = utils.TranslateDbError(err)
	if err == models.NotFoundError {
		log.Success(message)
		return nil, err
	}
	if err != models.AlreadyExistsError {
		log.Error(message, err)
		return nil, err
	}

	return s.ForumGet(f.Slug)
}

//func (s *Repository) ForumThreadCreate(thread *models.Thread) (*models.Thread, error) {
//
//}
//func (s *Repository) ForumUsers(slug string, so *models.SortOptions) ([]*models.User, error) {
//
//}

//func (s *Repository) UserCreate(u *models.User) ([]*models.User, error) {
//	message := logMessage + "UserCreate:"
//	log.Debug(message + "started")
//	user := jsonToDbModel(u)
//	query := userCreate
//	_, err := s.db.Queryx(query, user.Nickname, user.Fullname, user.About, user.Email)
//	if err == nil {
//		log.Success(message)
//		return nil, nil
//	}
//	err = utils.TranslateDbError(err)
//	if err != models.AlreadyExistsError {
//		return nil, err
//	}
//
//	var users []User
//	err = s.db.Select(&users, selectAlreadyExist, user.Nickname, user.Email)
//	if err != nil {
//		log.Error(message, err)
//		return nil, err
//	}
//	//if len(users) == 0 {
//	//	err = models.UnexpectedDbBehavior
//	//	log.Error(message, err)
//	//	return nil, err
//	//}
//	var us []*models.User
//	for _, u := range users {
//		us = append(us, dbToJsonModel(&u))
//	}
//	log.Success(message)
//	return us, models.AlreadyExistsError
//}
//
//func (s *Repository) UserProfileGet(nickname string) (*models.User, error) {
//	message := logMessage + "UserProfileGet:"
//	log.Debug(message + "started")
//	query := userProfileGet
//	user := new(User)
//	err := s.db.Get(user, query, nickname)
//	log.Error(message+"error = ", err)
//	if err == nil {
//		log.S(message + "[SUCCESS]")
//		return dbToJsonModel(user), nil
//	}
//	err = utils.TranslateDbError(err)
//	return nil, err
//}
//
//func (s *Repository) UserProfileUpdate(u *models.User) error {
//	message := logMessage + "UserProfileUpdate:"
//	log.Debug(message + "started")
//	user := jsonToDbModel(u)
//	query := userProfileUpdate
//	id := -1
//	err := s.db.Get(&id, query, user.Nickname, user.Fullname, user.About, user.Email)
//	if err == nil {
//		log.Debug(message + "[SUCCESS]")
//		//return nil, nil
//		return nil
//	}
//	err = utils.TranslateDbError(err)
//	return err
//}
