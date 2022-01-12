package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
)

const (
	logMessage   = "service:service:repository:"
	statusGet    = `select * from "Service" where id = 1`
	serviceClear = `truncate tpdb."User",  tpdb."Post", tpdb."Thread", tpdb."Forum" cascade;
		update "Service" set users = 0, posts = 0, threads = 0, forums = 0 where id = 1"`
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
	message := logMessage + "ServiceStatus:"
	log.Debug(message + "started")
	query := serviceClear
	_, err := s.db.Queryx(query)
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
	query := statusGet
	status := new(Status)
	err := s.db.Get(status, query)
	if err == nil {
		log.Success(message)
		return DbToJsonModel(status), nil
	}
	err = utils.TranslateDbError(err)
	log.Error(message, err)
	return nil, err
}
