package store

import (
	"database/sql"
	_ "github.com/lib/pq"
	"errors"
	"fmt"
	"context"
	"cc-first-project/advertisement-service/models"
)

func NewStore(db *sql.DB) Store {
	return &PostgresStore{Db: db}
}

type Store interface {
	SetCategory(adId string, category string, state string)error
}
type PostgresStore struct {
	Db *sql.DB
}

func (s *PostgresStore) SetCategory(adId, category, state string) error {
	_, err := s.Db.Query("UPDATE ads set category = $1, state = $2 where id = $3", category, state, adId)
	return err
}
