package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
)

const (
	logMessage   = "service:service:repository:"
	statusUser   = `select count(*) from "User";`
	statusPost   = `select count(*) from Post;`
	statusForum  = `select count(*) from Forum;`
	statusThread = `select count(*) from Thread;`
	serviceClear = `truncate "User",  Post, Thread, 
		Forum, Vote, ForumUsers cascade;`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) ServiceClear() error {
	message := logMessage + "ServiceClear:"
	log.Debug(message + "started")
	query := serviceClear
	rows, err := s.db.Queryx(query)
	if rows != nil {
		e := rows.Close()
		if e != nil {
			log.Error(message, e)
		}
	}
	if err == nil {
		log.Success(message)
		return nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return err
}

func (s *Repository) ServiceStatus() (*models.Status, error) {
	message := logMessage + "ServiceStatus:"
	log.Debug(message + "started")
	status := new(Status)
	err := s.db.Get(&status.User, statusUser)
	if err != nil {
		err = utils.TranslateDbError(err)
		log.Error(message, err)
		return nil, err
	}
	err = s.db.Get(&status.Forum, statusForum)
	if err != nil {
		err = utils.TranslateDbError(err)
		log.Error(message, err)
		return nil, err
	}
	err = s.db.Get(&status.Thread, statusThread)
	if err != nil {
		err = utils.TranslateDbError(err)
		log.Error(message, err)
		return nil, err
	}
	err = s.db.Get(&status.Post, statusPost)
	if err != nil {
		err = utils.TranslateDbError(err)
		log.Error(message, err)
		return nil, err
	}

	log.Success(message)
	return DbToJsonModel(status), err
}
