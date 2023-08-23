package storage

import (
	"database/sql"
	"os"

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
	dbName := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_DB_USERNAME")
	pwd := os.Getenv("POSTGRES_DB_PASSWORD")

	connStr := "user=" + user + " dbname=" + dbName + " password=" + pwd + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil

}

func (s *PostgresStore) Init() error {

	return s.createAccountTable()

}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
			id UUID PRIMARY KEY,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			number INTEGER UNIQUE,
			balance FLOAT,
			created_at TIMESTAMP
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) GetAccountByID(uuid.UUID) (*types.Account, error) {

	return nil, nil

}

func (s *PostgresStore) CreateAccount(account *types.Account) error {
	query := `INSERT INTO accounts
	(id, first_name, last_name, number, balance, created_at) 
	values
	($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query,
		account.ID,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*types.Account) error {

	return nil

}

func (s *PostgresStore) DeleteAccount(uuid.UUID) error {

	return nil

}
