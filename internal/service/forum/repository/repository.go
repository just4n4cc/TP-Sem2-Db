package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	userRepository "github.com/just4n4cc/tp-sem2-db/internal/service/user/repository"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
)

const (
	logMessage         = "service:forum:repository:"
	selectAlreadyExist = `select * from Forum
		where slug = $1`
	forumCreate = `insert into Forum
		(title, "user", slug)
		values($1, $2, $3) returning *`
	forumUsersNil = `select u.id, u.nickname, u.fullname, u.about, u.email 
		from "User"as u
		inner join ForumUsers as fu 
		on u.nickname = fu."user"
		where fu.forum = $1`
	forumUsers = `select u.id, u.nickname, u.fullname, u.about, u.email 
		from "User"as u
		inner join ForumUsers as fu 
		on u.nickname = fu."user"
		where fu.forum = $1
		order by "user" limit $2`
	forumUsersDesc = `select u.id, u.nickname, u.fullname, u.about, u.email 
		from "User"as u
		inner join ForumUsers as fu 
		on u.nickname = fu."user"
		where fu.forum = $1
		order by "user" desc limit $2`
	forumUsersSince = `select u.id, u.nickname, u.fullname, u.about, u.email 
		from "User"as u
		inner join ForumUsers as fu 
		on u.nickname = fu."user"
		where fu.forum = $1 and fu."user" > $3
		order by "user" limit $2`
	forumUsersSinceDesc = `select u.id, u.nickname, u.fullname, u.about, u.email 
		from "User"as u
		inner join ForumUsers as fu 
		on u.nickname = fu."user"
		where fu.forum = $1 and fu."user" < $3
		order by "user" desc limit $2`
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
	err := s.db.Get(forum, selectAlreadyExist, slug)
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
	err := s.db.Get(forum, query, forum.Title, forum.User, forum.Slug)
	if err == nil {
		log.Success(message)
		return dbToJsonModel(forum), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	if err == models.NotFoundError {
		log.Success(message)
		return nil, err
	}
	if err != models.AlreadyExistsError {
		log.Error(message, err)
		return nil, err
	}

	f, err = s.ForumGet(f.Slug)
	if err == nil {
		log.Success(message)
		return f, models.AlreadyExistsError
	}

	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}

func (s *Repository) ForumUsers(slug string, so *models.SortOptions) ([]*models.User, error) {
	message := logMessage + "ForumUsers:"
	log.Debug(message + "started")
	query := ""
	var args []interface{}
	args = append(args, slug)
	if so == nil {
		query = forumUsersNil
	} else {
		args = append(args, so.Limit)
		if so.Since == "" {
			if so.Desc {
				query += forumUsersDesc
			} else {
				query += forumUsers
			}
		} else {
			args = append(args, so.Since)
			if so.Desc {
				query += forumUsersSinceDesc
			} else {
				query += forumUsersSince
			}
		}
	}
	var users []userRepository.User
	err := s.db.Select(&users, query, args...)
	if err == nil {
		log.Success(message)
		return userRepository.DbArrayToJsonModel(users), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
