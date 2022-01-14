package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	"strconv"
	"strings"
)

const (
	logMessage   = "service:thread:repository:"
	threadBySlug = `select * from Thread
		where slug = $1`
	threadById = `select * from Thread
		where id = $1`
	threadCreate = `insert into Thread
		(title, author, forum, message, slug, created)
		values($1, $2, $3, $4, $5, $6)
		returning *`
	threadsByForumNil = `select * from Thread
		where forum = $1`
	threadsByForum = `select * from Thread
		where forum = $1 order by created limit $2`
	threadsByForumDesc = `select * from Thread
		where forum = $1 order by created desc limit $2`
	threadsByForumSince = `select * from Thread
		where forum = $1 and created >= $3 order by created limit $2`
	threadsByForumSinceDesc = `select * from Thread
		where forum = $1 and created <= $3 order by created desc limit $2`
	threadUpdateById = `update Thread
		set `
	threadUpdateByIdEnd = ` where id = $1
		returning *`
	threadUpdateBySlug = `update Thread
		set `
	threadUpdateBySlugEnd = ` where slug = $1
		returning *`
	threadVote = `insert into Vote
		(vote, threadid, "user")
		values ($1, $2, $3)`
	threadUpdateVote = `update Thread
		set votes = $2
		where id = $1`
	threadVoteGet = `
	select vote from Vote 
	where "user" = $2 and threadid = $1;`
	threadVoteUpdate = `
	update Vote
		set vote = $3
		where "user" = $2 and threadid = $1`
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
	//message := logMessage + "ThreadBySlug:"
	//log.Debug(message + "started")
	query := threadBySlug
	thread := new(Thread)
	err := s.db.Get(thread, query, slug)
	if err == nil {
		//log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	//log.Error(message, err)
	return nil, err
}

func (s *Repository) ThreadsByForum(forum string, so *models.SortOptions) ([]*models.Thread, error) {
	//message := logMessage + "ThreadsByForum:"
	//log.Debug(message + "started")
	query := ""
	var args []interface{}
	args = append(args, forum)
	if so == nil {
		query = threadsByForumNil
	} else {
		args = append(args, so.Limit)
		if so.Since == "" {
			if so.Desc {
				query = threadsByForumDesc
			} else {
				query = threadsByForum
			}
		} else {
			args = append(args, so.Since)
			if so.Desc {
				query = threadsByForumSinceDesc
			} else {
				query = threadsByForumSince
			}
		}
	}

	threads := new([]Thread)
	err := s.db.Select(threads, query, args...)
	if err == nil {
		var ts []*models.Thread
		for _, t := range *threads {
			ts = append(ts, dbToJsonModel(&t))
		}
		//log.Success(message)
		return ts, nil
	}
	err = utils.TranslateDbError(err)
	if err == models.NotFoundError {
		//log.Success(message)
	} else {
		//log.Error(message, err)
	}
	return nil, err
}

func (s *Repository) ThreadCreate(t *models.Thread) (*models.Thread, error) {
	//message := logMessage + "ThreadCreate:"
	//log.Debug(message + "started")
	thread := jsonToDbModel(t)
	query := threadCreate
	//log.Debug(*thread)

	err := s.db.Get(thread, query, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created)
	//log.Error(message, err)
	if err == nil {
		if t.Slug == "" {
			thread.Slug = ""
		}
		//log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	if err != models.AlreadyExistsError {
		if err == models.NotFoundError {
			//log.Success(message)
		} else {
			//log.Error(message, err)
		}
		return nil, err
	}

	t, err = s.ThreadBySlug(thread.Slug)
	if err == nil {
		//log.Success(message)
		return t, models.AlreadyExistsError
	}
	err = utils.TranslateDbError(err)
	//log.Error(message, err)
	return nil, err
}
func (s *Repository) ThreadById(id int32) (*models.Thread, error) {
	//message := logMessage + "ThreadById:"
	//log.Debug(message + "started")
	query := threadById
	thread := new(Thread)
	err := s.db.Get(thread, query, id)
	if err == nil {
		//log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	//log.Error(message, err)
	return nil, err
}
func (s *Repository) ThreadUpdateById(t *models.Thread) (*models.Thread, error) {
	//message := logMessage + "ThreadUpdateById:"
	//log.Debug(message + "started")
	query := threadUpdateById
	num := 2
	var args []interface{}
	args = append(args, t.Id)
	if t.Title != "" {
		query += "title = $" + strconv.Itoa(num) + ", "
		num++
		args = append(args, t.Title)
	}
	if t.Message != "" {
		query += "message = $" + strconv.Itoa(num) + ", "
		args = append(args, t.Message)
	}
	query = strings.TrimSuffix(query, ", ")
	query += threadUpdateByIdEnd
	thread := jsonToDbModel(t)
	err := s.db.Get(thread, query, args...)
	if err == nil {
		//log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	//log.Error(message, err)
	return nil, err
}
func (s *Repository) ThreadUpdateBySlug(t *models.Thread) (*models.Thread, error) {
	//message := logMessage + "ThreadUpdateBySlug:"
	//log.Debug(message + "started")
	query := threadUpdateBySlug
	num := 2
	var args []interface{}
	args = append(args, t.Slug)
	if t.Title != "" {
		query += "title = $" + strconv.Itoa(num) + ", "
		num++
		args = append(args, t.Title)
	}
	if t.Message != "" {
		query += "message = $" + strconv.Itoa(num) + ", "
		args = append(args, t.Message)
	}
	query = strings.TrimSuffix(query, ", ")
	query += threadUpdateBySlugEnd
	thread := jsonToDbModel(t)
	err := s.db.Get(thread, query, args...)
	if err == nil {
		//log.Success(message)
		return dbToJsonModel(thread), nil
	}
	err = utils.TranslateDbError(err)
	//log.Error(message, err)
	return nil, err
}
func (s *Repository) ThreadVote(t *models.Thread, v *models.Vote) (*models.Thread, error) {
	//message := logMessage + "ThreadVote:"
	//log.Debug(message + "started")
	query := threadVote
	rows, err := s.db.Queryx(query, v.Voice, t.Id, v.Nickname)
	if rows != nil {
		e := rows.Close()
		if e != nil {
			//log.Error(message, e)
		}
	}
	if err == nil { // INSERTED
		query = threadUpdateVote
		t.Votes += v.Voice
		rows, err = s.db.Queryx(query, t.Id, t.Votes)
		if rows != nil {
			e := rows.Close()
			if e != nil {
				//log.Error(message, e)
			}
		}
		if err == nil {
			//log.Success(message)
			return t, nil
		}
		err = utils.TranslateDbError(err)
		//log.Error(message, err)
		return nil, err
	}
	err = utils.TranslateDbError(err)
	if err == models.NotFoundError || err != models.AlreadyExistsError {
		//log.Error(message, err)
		return nil, err
	}
	query = threadVoteGet
	oldVote := int32(0)
	err = s.db.Get(&oldVote, query, t.Id, v.Nickname)
	if err != nil {
		err = utils.TranslateDbError(err)
		//log.Error(message, err)
		return nil, err
	}

	query = threadVoteUpdate
	rows, err = s.db.Queryx(query, t.Id, v.Nickname, v.Voice)
	if rows != nil {
		e := rows.Close()
		if e != nil {
			//log.Error(message, e)
		}
	}
	if err != nil {
		err = utils.TranslateDbError(err)
		//log.Error(message, err)
		return nil, err
	}

	if oldVote == v.Voice {
		//log.Success(message)
		return t, nil
	}
	t.Votes += v.Voice * 2
	query = threadUpdateVote
	rows, err = s.db.Queryx(query, t.Id, t.Votes)
	if rows != nil {
		e := rows.Close()
		if e != nil {
			//log.Error(message, e)
		}
	}
	if err == nil {
		//log.Success(message)
		return t, nil
	}

	err = utils.TranslateDbError(err)
	//log.Error(message, err)
	return nil, err
}
