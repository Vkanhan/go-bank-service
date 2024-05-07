package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Storage interface defines methods for interacting with a storage system.
type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
}

// PostgresStore is a struct representing a PostgreSQL database store.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new instance of PostgresStore and initializes the database connection.
func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=Kanhan@6520 sslmode=disable"
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

// Init initializes the PostgreSQL database by creating the account table.
func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

// CreateAccountTable creates the account table if it does not exist already.
func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`
	
	_, err := s.db.Exec(query)
	return err
}


// CreateAccount creates a new account in the PostgreSQL database.
func (s *PostgresStore) CreateAccount(*Account) error {
	return nil
}

// UpdateAccount updates an existing account in the PostgreSQL database.
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

// DeleteAccount deletes an account from the PostgreSQL database based on the given ID.
func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

// GetAccountByID retrieves an account from the PostgreSQL database based on the given ID.
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}

