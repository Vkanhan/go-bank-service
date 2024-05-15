package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Storage interface defines methods for interacting with a storage system.
type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
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
	if err := s.CreateAccountTable(); err != nil {
		return err
	}
	return nil
}

// CreateAccountTable creates the account table if it does not exist already.
func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name  VARCHAR(50) NOT NULL,
		number SERIAL,
		encrypted_password VARCHAR(100),
		balance SERIAL,
		created_at TIMESTAMP NOT NULL
	)`

	_, err := s.db.Exec(query)
	return err
}

// CreateAccount creates a new account in the PostgreSQL database.
func (s *PostgresStore) CreateAccount(acc *Account) error {

	query := `INSERT INTO account 
	(first_name, last_name, number, encrypted_password, balance, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateAccount updates an existing account in the PostgreSQL database.
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

// DeleteAccount deletes an account from the PostgreSQL database based on the given ID.
func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
	return err
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE number = $1", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return s.scanIntoAccounts(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", number)
}

// GetAccountByID retrieves an account from the PostgreSQL database based on the given ID.
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return s.scanIntoAccounts(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	// Query the database to get rows of accounts
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the retrieved accounts
	accounts := []*Account{}

	// Iterate over the rows
	for rows.Next() {
		// Corrected function invocation
		account, err := s.scanIntoAccounts(rows)
		if err != nil {
			return nil, err
		}

		// Append the scanned account to the accounts slice
		accounts = append(accounts, account)
	}
	// Return the slice of retrieved accounts
	return accounts, nil
}

func (s *PostgresStore) scanIntoAccounts(rows *sql.Rows) (*Account, error) {
	// Create a new Account object for each row
	account := new(Account)

	// Scan the values from the current row into the Account object fields
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
