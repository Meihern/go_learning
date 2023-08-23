package storage

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/meihern/go_learning/types"
)

type Storage interface {
	GetAccountByID(uuid.UUID) (*types.Account, error)
	CreateAccount(*types.Account) error
	UpdateAccount(*types.Account) error
	DeleteAccount(uuid.UUID) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	
	return nil, nil

}