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
	logMessage = "service:post:repository:"
	postGet    = `select * from tpdb."Post"
		where id = $1`
	postsCreate = `insert into tpdb."Post"
		(parent, author, message, forum, thread, created)
		values`
	postsByThread = `select * from tpdb."Post"
		where thread = $1`
	postUpdate = `update tpdb."Post"
		set message = $2, isEdited = true
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

func (s *Repository) PostGet(id int64) (*models.Post, error) {
	message := logMessage + "PostGet:"
	log.Debug(message + "started")
	query := postGet
	post := new(Post)
	err := s.db.Get(post, query, id)
	if err == nil {
		log.Success(message)
		return DbToJsonModel(post), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
func (s *Repository) PostUpdate(p *models.Post) (*models.Post, error) {
	message := logMessage + "PostUpdate:"
	log.Debug(message + "started")
	query := postUpdate
	post := new(Post)
	err := s.db.Get(post, query, p.Id, p.Message)
	if err == nil {
		log.Success(message)
		return DbToJsonModel(post), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}

func (s *Repository) PostsCreate(ps []*models.Post) ([]*models.Post, error) {
	message := logMessage + "PostsCreate:"
	log.Debug(message + "started")
	query := postsCreate
	var args []interface{}
	num := 1
	for _, p := range ps {
		args = append(args, p.Parent)
		args = append(args, p.Author)
		args = append(args, p.Message)
		args = append(args, p.Created)
		query += " ($" + strconv.Itoa(num) + ", $"
		num++
		query += strconv.Itoa(num) + ", $"
		num++
		query += strconv.Itoa(num) + ", '" + p.Forum + "', " + strconv.Itoa(int(p.Thread)) + ", $"
		num++
		query += strconv.Itoa(num) + "),"
		num++
	}
	query = strings.TrimSuffix(query, ",")
	query = query + " returning *"
	log.Debug(message + "query = " + query)
	var posts []Post
	err := s.db.Select(&posts, query, args...)
	if err == nil {
		for i, p := range posts {
			ps[i] = DbToJsonModel(&p)
		}
		log.Success(message)
		return ps, nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}

func (s *Repository) PostsByThread(id int32, so *models.SortOptions) ([]*models.Post, error) {
	message := logMessage + "PostsByThread:"
	log.Debug(message + "started")
	query := postsByThread + utils.SortOptionsToSubquery(so, "created")
	var posts []Post
	err := s.db.Select(&posts, query, id)
	if err == nil {
		var ps []*models.Post
		for _, p := range posts {
			ps = append(ps, DbToJsonModel(&p))
		}
		log.Success(message)
		return ps, nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
