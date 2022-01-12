package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
)

const (
	logMessage   = "service:thread:repository:"
	threadBySlug = `select * from tpdb."Thread"
		where slug = $1`
	threadById = `select * from tpdb."Thread"
		where id = $1`
	threadCreate = `insert into tpdb."Thread"
		(title, author, forum, message, slug, created)
		values($1, $2, $3, $4, $5, $6)
		returning id`
	threadsByForum = `select * from tpdb."Thread"
		where forum = $1`
	threadUpdateById = `update tpdb."Thread"
		set title = $2, message = $3
		where id = $1
		returning *`
	threadUpdateBySlug = `update tpdb."Thread"
		set title = $2, message = $3
		where slug = $1
		returning *`
	threadVoteById = `update tpdb."Thread"
		set votes = votes + $2
		where id = $1
		returning *`
	threadVoteBySlug = `update tpdb."Thread"
		set votes = votes + $2
		where slug = $1
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

func (s *Repository) ThreadBySlug(slug string) (*models.Thread, error) {
	message := logMessage + "ThreadBySlug:"
	log.Debug(message + "started")
	query := threadBySlug
	thread := new(Thread)
	err := s.db.Get(thread, query, slug)
	if err == nil {
		log.Success(message)
		return dbToJsonModel(thread), nil
	}
	log.Error(message, err)
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}

func (s *Repository) ThreadsByForum(forum string, so *models.SortOptions) ([]*models.Thread, error) {
	message := logMessage + "ThreadsByForum:"
	log.Debug(message + "started")
	query := threadsByForum
	query += utils.SortOptionsToSubquery(so, "created")
	//sign := ""
	//order := ""
	//if !so.Desc {
	//	sign = "<="
	//	order = "desc"
	//} else {
	//	sign = ">="
	//}
	//if so.Since != "" {
	//	query += " and created " + sign
	//}
	//query += " order by created " + order + " limit " + string(so.Limit)
	threads := new([]Thread)
	err := s.db.Select(threads, query, forum)
	if err == nil {
		var ts []*models.Thread
		for _, t := range *threads {
			ts = append(ts, dbToJsonModel(&t))
		}
		log.Success(message)
		return ts, nil
	}
	err = utils.TranslateDbError(err)
	if err == models.NotFoundError {
		log.Success(message)
	} else {
		log.Error(message, err)
	}
	return nil, err
}

func (s *Repository) ThreadCreate(t *models.Thread) (*models.Thread, error) {
	message := logMessage + "ThreadCreate:"
	log.Debug(message + "started")
	thread := jsonToDbModel(t)
	query := threadCreate
	var id int32
	log.Debug(*thread)
	err := s.db.Get(&id, query, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created)
	log.Error(message, err)
	if err == nil {
		thread.Id = id
		log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	if err != models.AlreadyExistsError {
		if err == models.NotFoundError {
			log.Success(message)
		} else {
			log.Error(message, err)
		}
		return nil, err
	}

	t, err = s.ThreadBySlug(thread.Slug)
	if err == nil {
		log.Success(message)
		return t, models.AlreadyExistsError
	}
	log.Error(message, err)
	return nil, err
}
func (s *Repository) ThreadById(id int32) (*models.Thread, error) {
	message := logMessage + "ThreadByrId:"
	log.Debug(message + "started")
	query := threadById
	thread := new(Thread)
	err := s.db.Get(thread, query, id)
	if err == nil {
		log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
func (s *Repository) ThreadUpdateById(t *models.Thread) (*models.Thread, error) {
	message := logMessage + "ThreadById:"
	log.Debug(message + "started")
	query := threadUpdateById
	thread := jsonToDbModel(t)
	err := s.db.Get(thread, query, t.Id, t.Title, t.Message)
	if err == nil {
		log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
func (s *Repository) ThreadUpdateBySlug(t *models.Thread) (*models.Thread, error) {
	message := logMessage + "ThreadBySlug:"
	log.Debug(message + "started")
	query := threadUpdateBySlug
	thread := jsonToDbModel(t)
	err := s.db.Get(thread, query, t.Slug, t.Title, t.Message)
	if err == nil {
		log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
func (s *Repository) ThreadVoteById(id int32, vote *models.Vote) (*models.Thread, error) {
	message := logMessage + "ThreadVoteById:"
	log.Debug(message + "started")
	query := threadVoteById
	thread := new(Thread)
	err := s.db.Get(thread, query, id, vote.Voice)
	if err == nil {
		log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
func (s *Repository) ThreadVoteBySlug(slug string, vote *models.Vote) (*models.Thread, error) {
	message := logMessage + "ThreadVoteBySlug:"
	log.Debug(message + "started")
	query := threadVoteBySlug
	thread := new(Thread)
	err := s.db.Get(thread, query, slug, vote.Voice)
	if err == nil {
		log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
