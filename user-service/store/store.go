package store

import (
	"database/sql"
	_ "github.com/lib/pq"
	"errors"
	"fmt"
	"context"
	"cc-first-project/user-service/models"
)

func NewStore(db *sql.DB) Store {
	_, err := db.Query("create table if not exists ads (id varchar(255), description varchar(255), email varchar(255), state varchar(255), category varchar(255), image varchar(255))")
	if err != nil {
		panic(err)
		return nil
	}
	return &PostgresStore{Db: db}
}

type Store interface {
	CreateAdvertisement(ctx context.Context, ad *models.Advertisement) error
	GetAdvertisement(ctx context.Context, adId string) (*models.Advertisement, error)
}
type PostgresStore struct {
	Db *sql.DB
}

func (s *PostgresStore) CreateAdvertisement(ctx context.Context, ad *models.Advertisement) error {
	_, err := s.Db.Query("INSERT INTO ads (id, description, email, state, category, image) VALUES ($1, $2, $3, $4, $5, $6)", ad.Id, ad.Description, ad.Email, ad.State, ad.Category, ad.Image)
	return err
}

func (s *PostgresStore) GetAdvertisement(ctx context.Context, adId string) (*models.Advertisement, error) {
	rows, err := s.Db.Query("SELECT * FROM ads WHERE id = $1", adId)
	if err != nil {
		return nil, errors.New("error while getting advertisement" + err.Error())
	}

	for rows.Next() {
		var result models.Advertisement
		err = rows.Scan(&result.Id, &result.Description, &result.Email, &result.State, &result.Category, &result.Image)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Version: %s\n", result)
		return &result, nil
	}
	return nil, errors.New("not found")
}
