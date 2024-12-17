package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	// "wallet-service/internal/model"

	"github.com/google/uuid"
	// "google.golang.org/grpc"
)

// Wallet represents a wallet model.
type Wallet struct {
	ID       uuid.UUID
	Balance  int64
	Currency string
}

// User represents a user model for the login system.
type User struct {
	ID       uuid.UUID
	Username string
	Password string // Stored as a hash
}

// WalletRepository handles wallet-related database operations.
type WalletRepository struct {
	db *sql.DB
	mu sync.Mutex // To handle concurrent operations
}

// WalletRepositoryInterface defines the contract for wallet operations.
type WalletRepositoryInterface interface {
	GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error)
	UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error
	CreateUser(ctx context.Context, username, password string) (uuid.UUID, error)
	Login(ctx context.Context, username, password string) (*User, error)
}

// Config holds database configuration details.
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfigFromEnv loads database configuration from environment variables.
func LoadConfigFromEnv() Config {
	return Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

// NewPostgresDB initializes a new PostgreSQL database connection.
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// Connection check
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// Connection pool settings
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

// NewWalletRepository creates a new WalletRepository instance.
func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// CreateUser creates a new user in the database.
func (r *WalletRepository) CreateUser(ctx context.Context, username, password string) (uuid.UUID, error) {
	hashedPassword := hashPassword(password)
	userID := uuid.New()
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", userID, username, hashedPassword)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

// Login validates user credentials and retrieves user details.
func (r *WalletRepository) Login(ctx context.Context, username, password string) (*User, error) {
	var user User
	var hashedPassword string

	err := r.db.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	if hashPassword(password) != hashedPassword {
		return nil, errors.New("invalid username or password")
	}
	return &user, nil
}

// GetBalance retrieves the balance of a specific wallet.
func (r *WalletRepository) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	var balance int64
	err := r.db.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE id = $1", walletID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, errors.New("wallet not found")
	}
	return balance,
		err
}

// UpdateBalance updates the wallet balance after acquiring a lock.
func (r *WalletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var balance int64
	err = tx.QueryRow("SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", walletID).Scan(&balance)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return errors.New("wallet not found")
	} else if err != nil {
		tx.Rollback()
		return err
	}

	newBalance := balance + amount
	if newBalance < 0 {
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	_, err = tx.Exec("UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3", newBalance, time.Now(), walletID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// hashPassword hashes a password using SHA-256.
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
