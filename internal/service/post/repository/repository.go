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
	postGet    = `select * from Post
		where id = $1`
	postsCreate = `insert into Post
		(parent, author, message, forum, thread, created)
		values`

	postsByThreadNil = `select * from Post
		where thread = $1`
	postsByThreadFlat = `select * from Post
		where thread = $1 order by created, id limit $2`
	postsByThreadFlatDesc = `select * from Post
		where thread = $1 order by created desc, id desc limit $2`
	postsByThreadFlatSince = `select * from Post
		where thread = $1 and id > $3 order by created, id limit $2`
	postsByThreadFlatSinceDesc = `select * from Post
		where thread = $1 and id < $3 order by created desc, id desc limit $2`

	postsByThreadTree = `select * from Post
		where thread = $1 order by path, id limit $2`
	postsByThreadTreeDesc = `select * from Post
		where thread = $1 order by path desc, id desc limit $2`
	postsByThreadTreeSince = `select * from Post
		where thread = $1 and path > (select path from Post where id = $3) order by path, id limit $2`
	postsByThreadTreeSinceDesc = `select * from Post
		where thread = $1 and path < (select path from Post where id = $3) order by path desc, id desc limit $2`

	postsByThreadPTree = `select * from Post
		where path[1] in (select id from Post 
							where thread = $1 and parent = 0 
							order by id limit $2)
		order by path, id`
	postsByThreadPTreeDesc = `select * from Post
		where path[1] in (select id from Post 
							where thread = $1 and parent = 0 
							order by id desc limit $2)
		order by path[1] desc, path, id`
	postsByThreadPTreeSince = `select * from Post
		where path[1] in (select id from Post 
							where thread = $1 and parent = 0 and path[1] > ( select path[1] from Post
																				where id = $3
																			)
							order by id limit $2)
		order by path, id`
	postsByThreadPTreeSinceDesc = `select * from Post
		where path[1] in (select id from Post 
							where thread = $1 and parent = 0 and path[1] < ( select path[1] from Post
																				where id = $3
																			)
							order by id desc limit $2)
		order by path[1] desc, path, id`

	postUpdate = `update Post
		set message = $2
		where id = $1
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
	//message := logMessage + "PostGet:"
	//log.Debug(message + "started")
	query := postGet
	post := new(Post)
	err := s.db.Get(post, query, id)
	if err == nil {
		//log.Success(message)
		return DbToJsonModel(post), nil
	}
	err = utils.TranslateDbError(err)
	//log.Error(message, err)
	return nil, err
}
func (s *Repository) PostUpdate(p *models.Post) (*models.Post, error) {
	//message := logMessage + "PostUpdate:"
	//log.Debug(message + "started")
	query := postUpdate
	post := new(Post)
	err := s.db.Get(post, query, p.Id, p.Message)
	if err == nil {
		//log.Success(message)
		return DbToJsonModel(post), nil
	}
	err = utils.TranslateDbError(err)
	//log.Error(message, err)
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
	var posts []Post
	err := s.db.Select(&posts, query, args...)
	log.Error(message, err)
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
	//message := logMessage + "PostsByThread:"
	//log.Debug(message + "started")
	query := ""
	var args []interface{}
	args = append(args, id)
	if so == nil {
		query = postsByThreadNil
	} else {
		args = append(args, so.Limit)
		if so.Sort == "flat" {
			if so.Since == "" {
				if so.Desc {
					query = postsByThreadFlatDesc
				} else {
					query = postsByThreadFlat
				}
			} else {
				args = append(args, so.Since)
				if so.Desc {
					query = postsByThreadFlatSinceDesc
				} else {
					query = postsByThreadFlatSince
				}
			}
		} else if so.Sort == "tree" {
			if so.Since == "" {
				if so.Desc {
					query = postsByThreadTreeDesc
				} else {
					query = postsByThreadTree
				}
			} else {
				args = append(args, so.Since)
				if so.Desc {
					query = postsByThreadTreeSinceDesc
				} else {
					query = postsByThreadTreeSince
				}
			}
		} else {
			if so.Since == "" {
				if so.Desc {
					query = postsByThreadPTreeDesc
				} else {
					query = postsByThreadPTree
				}
			} else {
				args = append(args, so.Since)
				if so.Desc {
					query = postsByThreadPTreeSinceDesc
				} else {
					query = postsByThreadPTreeSince
				}
			}
		}
	}

	var posts []Post
	err := s.db.Select(&posts, query, args...)
	if err == nil {
		var ps []*models.Post
		for _, p := range posts {
			ps = append(ps, DbToJsonModel(&p))
		}
		//log.Success(message)
		return ps, nil
	}
	err = utils.TranslateDbError(err)
	//log.Error(message, err)
	return nil, err
}
