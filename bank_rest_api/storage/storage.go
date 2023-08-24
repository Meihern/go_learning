package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/meihern/go_learning/types"
)

type Storage interface {
	GetAccounts() ([]*types.Account, error)
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
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`
	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	query = `ALTER TABLE IF EXISTS accounts
			ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP`
	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {
	query := `SELECT * FROM accounts`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil

}

func (s *PostgresStore) GetAccountByID(id uuid.UUID) (*types.Account, error) {

	query := `SELECT * FROM accounts WHERE id = $1`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account %s not found", id)

}

func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	account := new(types.Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return account, nil

}

func (s *PostgresStore) CreateAccount(account *types.Account) error {
	query := `INSERT INTO accounts
	(id, first_name, last_name, number, balance, created_at, updated_at) 
	values
	($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.Query(
		query,
		account.ID,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(account *types.Account) error {
	query := `UPDATE accounts
	SET first_name = $2, last_name = $3, updated_at = $4
	WHERE id = $1`

	_, err := s.db.Query(query,
		account.ID,
		account.FirstName,
		account.LastName,
		account.UpdatedAt)
	if err != nil {
		return err
	}

	return nil

}

func (s *PostgresStore) DeleteAccount(id uuid.UUID) error {
	query := `DELETE FROM accounts where id = $1`

	_, err := s.db.Query(query, id)
	if err != nil {
		return fmt.Errorf("Account %s not found", id)
	}

	return nil
}
