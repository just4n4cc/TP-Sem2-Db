package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/just4n4cc/tp-sem2-db/internal/models"
)

type Repository struct {
	db *sql.DB
}

func (s *Repository) ThreadsByForum(slug string, so *models.SortOptions) ([]*models.Thread, error) {
	return nil, nil
}
