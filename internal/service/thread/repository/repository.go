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
		returning *`
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
	threadVote = `insert into tpdb."Vote"
		(vote, threadid, threadslug, "user")
		values ($1, $2, $3, $4)`
	threadUpdateVote = `update tpdb."Thread"
		set votes = $2
		where id = $1`
	threadVoteGet = `
	select vote from tpdb."Vote" 
	where threadid = $1 and "user" = $2;`
	threadVoteUpdate = `
	update tpdb."Vote"
		set vote = $3
		where threadid = $1 and "user" = $2`
	//threadVoteById = `update tpdb."Thread"
	//	set votes = votes + $2
	//	where id = $1
	//	returning *`
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
	log.Debug(message+"query = ", query)
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

func (s *Repository) ThreadCreate(t *models.Thread, slug string) (*models.Thread, error) {
	message := logMessage + "ThreadCreate:"
	log.Debug(message + "started")
	thread := jsonToDbModel(t)
	query := threadCreate
	log.Debug(*thread)
	err := s.db.Get(thread, query, thread.Title, thread.Author, thread.Forum, thread.Message, slug, thread.Created)
	log.Error(message, err)
	if err == nil {
		if t.Slug == "" {
			thread.Slug = ""
		}
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
func (s *Repository) ThreadVote(t *models.Thread, v *models.Vote) (*models.Thread, error) {
	message := logMessage + "ThreadVote:"
	log.Debug(message + "started")
	query := threadVote
	//if vote.Voice == -1 {
	//	vote.Voice = -2
	//}
	log.Debug(message+"vote = ", *v)
	log.Debug(message + " query = " + query)
	//err := s.db.Get(thread, query, id, vote.Voice)
	rows, err := s.db.Queryx(query, v.Voice, t.Id, t.Slug, v.Nickname)
	if rows != nil {
		e := rows.Close()
		if e != nil {
			log.Error(message, e)
		}
	}
	if err == nil { // INSERTED
		log.Debug(message + "INSERTED")
		log.Debug(message + " query = " + query)
		query = threadUpdateVote
		t.Votes += v.Voice
		rows, err = s.db.Queryx(query, t.Id, t.Votes)
		if rows != nil {
			e := rows.Close()
			if e != nil {
				log.Error(message, e)
			}
		}
		if err == nil {
			log.Success(message)
			return t, nil
		}
		err = utils.TranslateDbError(err)
		log.Error(message, err)
		return nil, err
	}
	err = utils.TranslateDbError(err)
	if err == models.NotFoundError || err != models.AlreadyExistsError {
		log.Error(message, err)
		return nil, err
	}
	query = threadVoteGet
	oldVote := int32(0)
	err = s.db.Get(&oldVote, query, t.Id, v.Nickname)
	if err != nil {
		err = utils.TranslateDbError(err)
		log.Error(message, err)
		return nil, err
	}

	query = threadVoteUpdate
	rows, err = s.db.Queryx(query, t.Id, v.Nickname, v.Voice)
	if rows != nil {
		e := rows.Close()
		if e != nil {
			log.Error(message, e)
		}
	}
	if err != nil {
		err = utils.TranslateDbError(err)
		log.Error(message, err)
		return nil, err
	}

	if oldVote == v.Voice {
		log.Success(message)
		return t, nil
	}
	t.Votes += v.Voice * 2
	query = threadUpdateVote
	rows, err = s.db.Queryx(query, t.Id, t.Votes)
	if rows != nil {
		e := rows.Close()
		if e != nil {
			log.Error(message, e)
		}
	}
	if err == nil {
		log.Success(message)
		return t, nil
	}

	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
