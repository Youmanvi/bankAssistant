package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Database holds the database connection
type Database struct {
	conn *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(config *Config) (*Database, error) {
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBSSLMode,
	)

	// Open database connection
	db, err := sql.Open(config.DBDriver, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to PostgreSQL database: %s", config.DBName)

	return &Database{conn: db}, nil
}

// Initialize creates all necessary tables
func (d *Database) Initialize() error {
	queries := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(20) UNIQUE NOT NULL,
			pin VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			address TEXT,
			date_of_birth DATE,
			ssn VARCHAR(20),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Sessions table
		`CREATE TABLE IF NOT EXISTS sessions (
			id SERIAL PRIMARY KEY,
			token VARCHAR(255) UNIQUE NOT NULL,
			user_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL
		)`,

		// Accounts table
		`CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			account_id VARCHAR(255) UNIQUE NOT NULL,
			user_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
			type VARCHAR(50) NOT NULL,
			balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Transactions table (for statements)
		`CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			account_id VARCHAR(255) NOT NULL REFERENCES accounts(account_id) ON DELETE CASCADE,
			from_account VARCHAR(255),
			to_account VARCHAR(255),
			amount DECIMAL(15, 2) NOT NULL,
			transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			description TEXT
		)`,

		// Create indexes
		`CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone)`,
		`CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_accounts_account_id ON accounts(account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id)`,
	}

	for _, query := range queries {
		if _, err := d.conn.Exec(query); err != nil {
			return fmt.Errorf("failed to execute initialization query: %w", err)
		}
	}

	log.Printf("Database schema initialized successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.conn.Close()
}

// ============================================================================
// User Operations
// ============================================================================

// CreateUser inserts a new user into the database
func (d *Database) CreateUser(user *AuthUser) error {
	query := `
		INSERT INTO users (user_id, phone, pin, name, email, address, date_of_birth, ssn)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := d.conn.Exec(query,
		user.UserID,
		user.Phone,
		user.PIN,
		user.Name,
		user.Email,
		user.Address,
		user.DateOfBirth,
		user.SSN,
	)

	return err
}

// GetUserByPhone retrieves a user by phone number
func (d *Database) GetUserByPhone(phone string) (*AuthUser, error) {
	query := `
		SELECT user_id, phone, pin, name, email, address, date_of_birth, ssn, created_at
		FROM users
		WHERE phone = $1
	`

	row := d.conn.QueryRow(query, phone)

	user := &AuthUser{
		Accounts: []string{},
	}

	err := row.Scan(
		&user.UserID,
		&user.Phone,
		&user.PIN,
		&user.Name,
		&user.Email,
		&user.Address,
		&user.DateOfBirth,
		&user.SSN,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	// Fetch user's accounts
	accounts, err := d.GetUserAccounts(user.UserID)
	if err == nil {
		user.Accounts = accounts
	}

	return user, nil
}

// GetUserByID retrieves a user by user ID
func (d *Database) GetUserByID(userID string) (*AuthUser, error) {
	query := `
		SELECT user_id, phone, pin, name, email, address, date_of_birth, ssn, created_at
		FROM users
		WHERE user_id = $1
	`

	row := d.conn.QueryRow(query, userID)

	user := &AuthUser{
		Accounts: []string{},
	}

	err := row.Scan(
		&user.UserID,
		&user.Phone,
		&user.PIN,
		&user.Name,
		&user.Email,
		&user.Address,
		&user.DateOfBirth,
		&user.SSN,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	// Fetch user's accounts
	accounts, err := d.GetUserAccounts(user.UserID)
	if err == nil {
		user.Accounts = accounts
	}

	return user, nil
}

// ListAllUsers returns all users
func (d *Database) ListAllUsers() ([]*AuthUser, error) {
	query := `
		SELECT user_id, phone, pin, name, email, address, date_of_birth, ssn, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*AuthUser

	for rows.Next() {
		user := &AuthUser{
			Accounts: []string{},
		}

		err := rows.Scan(
			&user.UserID,
			&user.Phone,
			&user.PIN,
			&user.Name,
			&user.Email,
			&user.Address,
			&user.DateOfBirth,
			&user.SSN,
			&user.CreatedAt,
		)

		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}

		// Fetch accounts for this user
		accounts, err := d.GetUserAccounts(user.UserID)
		if err == nil {
			user.Accounts = accounts
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

// ============================================================================
// Session Operations
// ============================================================================

// CreateSession inserts a new session into the database
func (d *Database) CreateSession(session *UserSession) error {
	query := `
		INSERT INTO sessions (token, user_id, created_at, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := d.conn.Exec(query,
		session.Token,
		session.UserID,
		time.Unix(session.CreatedAt, 0),
		time.Unix(session.ExpiresAt, 0),
	)

	return err
}

// GetSessionByToken retrieves a session by token
func (d *Database) GetSessionByToken(token string) (*UserSession, error) {
	query := `
		SELECT token, user_id, created_at, expires_at
		FROM sessions
		WHERE token = $1
	`

	row := d.conn.QueryRow(query, token)

	session := &UserSession{}
	var createdAt, expiresAt time.Time

	err := row.Scan(
		&session.Token,
		&session.UserID,
		&createdAt,
		&expiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, err
	}

	session.CreatedAt = createdAt.Unix()
	session.ExpiresAt = expiresAt.Unix()

	return session, nil
}

// DeleteSession removes a session by token
func (d *Database) DeleteSession(token string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	_, err := d.conn.Exec(query, token)
	return err
}

// DeleteExpiredSessions removes all expired sessions
func (d *Database) DeleteExpiredSessions() error {
	query := `DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := d.conn.Exec(query)
	return err
}

// ============================================================================
// Account Operations
// ============================================================================

// CreateAccount inserts a new account into the database
func (d *Database) CreateAccount(accountID, userID, accountType string, balance float64) error {
	query := `
		INSERT INTO accounts (account_id, user_id, type, balance)
		VALUES ($1, $2, $3, $4)
	`

	_, err := d.conn.Exec(query, accountID, userID, accountType, balance)
	return err
}

// GetUserAccounts retrieves all account IDs for a user
func (d *Database) GetUserAccounts(userID string) ([]string, error) {
	query := `
		SELECT account_id FROM accounts WHERE user_id = $1 ORDER BY created_at
	`

	rows, err := d.conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []string

	for rows.Next() {
		var accountID string
		if err := rows.Scan(&accountID); err != nil {
			log.Printf("Error scanning account: %v", err)
			continue
		}
		accounts = append(accounts, accountID)
	}

	return accounts, rows.Err()
}

// GetAccountDetails retrieves account information
func (d *Database) GetAccountDetails(accountID string) (map[string]interface{}, error) {
	query := `
		SELECT account_id, user_id, type, balance, created_at
		FROM accounts
		WHERE account_id = $1
	`

	row := d.conn.QueryRow(query, accountID)

	var accID, userID, accType string
	var balance float64
	var createdAt time.Time

	err := row.Scan(&accID, &userID, &accType, &balance, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"account_id": accID,
		"user_id":    userID,
		"type":       accType,
		"balance":    balance,
		"created_at": createdAt.Unix(),
	}, nil
}

// GetAccountBalance retrieves the balance for an account
func (d *Database) GetAccountBalance(accountID string) (float64, error) {
	query := `SELECT balance FROM accounts WHERE account_id = $1`

	row := d.conn.QueryRow(query, accountID)

	var balance float64
	err := row.Scan(&balance)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("account not found")
	}

	return balance, err
}

// UpdateAccountBalance updates the balance for an account
func (d *Database) UpdateAccountBalance(accountID string, newBalance float64) error {
	query := `
		UPDATE accounts
		SET balance = $1, updated_at = CURRENT_TIMESTAMP
		WHERE account_id = $2
	`

	_, err := d.conn.Exec(query, newBalance, accountID)
	return err
}

// ============================================================================
// Transaction Operations
// ============================================================================

// RecordTransaction records a transaction for an account
func (d *Database) RecordTransaction(accountID, fromAccount, toAccount, description string, amount float64) error {
	query := `
		INSERT INTO transactions (account_id, from_account, to_account, amount, description)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := d.conn.Exec(query, accountID, fromAccount, toAccount, amount, description)
	return err
}

// GetAccountTransactions retrieves transactions for an account
func (d *Database) GetAccountTransactions(accountID string, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, account_id, from_account, to_account, amount, transaction_date, description
		FROM transactions
		WHERE account_id = $1
		ORDER BY transaction_date DESC
		LIMIT $2
	`

	rows, err := d.conn.Query(query, accountID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []map[string]interface{}

	for rows.Next() {
		var id int
		var accID, fromAcc, toAcc, desc string
		var amount float64
		var txDate time.Time

		err := rows.Scan(&id, &accID, &fromAcc, &toAcc, &amount, &txDate, &desc)
		if err != nil {
			log.Printf("Error scanning transaction: %v", err)
			continue
		}

		transaction := map[string]interface{}{
			"id":              id,
			"account_id":      accID,
			"from_account":    fromAcc,
			"to_account":      toAcc,
			"amount":          amount,
			"transaction_date": txDate.Unix(),
			"description":     desc,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, rows.Err()
}
