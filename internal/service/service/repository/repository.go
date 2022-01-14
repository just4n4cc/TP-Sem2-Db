package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
)

const (
	logMessage = "service:service:repository:"
	//statusGet    = `select * from "Service" where id = 1`
	statusUser   = `select count(*) from tpdb."User";`
	statusPost   = `select count(*) from tpdb."Post";`
	statusForum  = `select count(*) from tpdb."Forum";`
	statusThread = `select count(*) from tpdb."Thread";`
	serviceClear = `truncate tpdb."User",  tpdb."Post", tpdb."Thread", 
		tpdb."Forum", tpdb."Vote", tpdb."ForumUsers" cascade;`
	//update "Service" set users = 0, posts = 0, threads = 0, forums = 0 where id = 1`
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
	//message := logMessage + "ServiceClear:"
	//log.Debug(message + "started")
	query := serviceClear
	rows, err := s.db.Queryx(query)
	if rows != nil {
		e := rows.Close()
		if e != nil {
			//log.Error(message, e)
		}
	}
	if err == nil {
		//log.Success(message)
		return nil
	}
	err = utils.TranslateDbError(err)
	//log.Error(message, err)
	return err
}

func (s *Repository) ServiceStatus() (*models.Status, error) {
	//message := logMessage + "ServiceStatus:"
	//log.Debug(message + "started")
	status := new(Status)
	err := s.db.Get(&status.User, statusUser)
	if err != nil {
		err = utils.TranslateDbError(err)
		//log.Error(message, err)
		return nil, err
	}
	err = s.db.Get(&status.Forum, statusForum)
	if err != nil {
		err = utils.TranslateDbError(err)
		//log.Error(message, err)
		return nil, err
	}
	err = s.db.Get(&status.Thread, statusThread)
	if err != nil {
		err = utils.TranslateDbError(err)
		//log.Error(message, err)
		return nil, err
	}
	err = s.db.Get(&status.Post, statusPost)
	if err != nil {
		err = utils.TranslateDbError(err)
		//log.Error(message, err)
		return nil, err
	}

	//log.Success(message)
	return DbToJsonModel(status), err
}
