package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	// "wallet-service/internal/model"

	// pb "wallet/internal/grpc/proto-exchange/grpc/pb"
	pb "github.com/SafetyDuck5676/grpc_duck/proto-exchange"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// Wallet represents a wallet model.
type Wallet struct {
	ID       uuid.UUID
	Balance  int64
	Currency string
}

// User represents a user model for the login system.
type User struct {
	ID       int32
	Username string
	Email    string
	Password string // Stored as a hash
}

type Token struct {
	Token string `json:"token"`
}

// WalletRepository handles wallet-related database operations.
type WalletRepository struct {
	db *sql.DB
	mu sync.Mutex // To handle concurrent operations
}

// WalletRepositoryInterface defines the contract for wallet operations.
type WalletRepositoryInterface interface {
	GetBalance(ctx context.Context, username string) (map[string]int32, error)
	UpdateBalance(ctx context.Context, uid int32, amount int32, currency string) error
	GetExchangeRates(ctx context.Context) (map[string]float64, error)
	GetExchangeRate(ctx context.Context, from string, to string) (float32, error)
	RegisterUser(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, username, password string) (Token, error)
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

// GetBalance retrieves the balance of a specific wallet.
func (r *WalletRepository) GetBalance(ctx context.Context, username string) (map[string]int32, error) {
	log.Println("Get balances")
	log.Println(username)
	var balances = make(map[string]int32)

	// Query to get the balance of a user's wallet
	query := "SELECT balance, currency FROM mydb.users JOIN mydb.wallets AS wallets ON users.id = wallets.user_id JOIN mydb.balances AS balance ON wallets.id = balance.wallet_id JOIN mydb.currencies ON mydb.currencies.id = balance.currency_id WHERE users.username = $1"
	rows, err := r.db.QueryContext(ctx, query, username)

	// Handle rows and scan the results
	for rows.Next() {
		var balance int32
		var currency string
		if err := rows.Scan(&balance, &currency); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		balances[currency] = balance
	}

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return balances, err
}

// UpdateBalance updates the wallet balance after acquiring a lock.
func (r *WalletRepository) UpdateBalance(ctx context.Context, uid int32, amount int32, currency string) error {
	log.Println("Updating balance")
	log.Println(uid, amount, currency)

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Select currency ID by currency name
	var currency_id int32
	err = tx.QueryRowContext(ctx, "SELECT id FROM mydb.currencies WHERE currency = $1", currency).Scan(&currency_id)
	if err == sql.ErrNoRows {
		log.Println("Currency not found")
		tx.Rollback()
		return errors.New("currency not found")
	} else if err != nil {
		log.Printf("Error with currency: %v", err)
		tx.Rollback()
		return err
	}

	// Select balance ID by user ID and currency ID
	var balance_id int32
	err = tx.QueryRowContext(ctx, "SELECT mydb.balances.id FROM mydb.wallets INNER JOIN mydb.balances ON mydb.balances.wallet_id = mydb.wallets.id  WHERE mydb.wallets.user_id = $1 AND mydb.balances.currency_id = $2", uid, currency_id).Scan(&balance_id)
	if err == sql.ErrNoRows {
		log.Println("Wallet not found")
		tx.Rollback()
		return errors.New("wallet not found")
	} else if err != nil {
		log.Printf("Error with wallet: %v", err)
		tx.Rollback()
		return err
	}

	// Update the balance
	var newBalance int32
	err = tx.QueryRowContext(ctx, "UPDATE mydb.balances SET balance = balance + $1 WHERE id = $2 RETURNING balance", amount, balance_id).Scan(&newBalance)
	if err != nil {
		log.Printf("Error with update: %v", err)
		log.Printf("Error with update: %v", "UPDATE mydb.balances SET balance = balance + $1 WHERE id = $2 RETURNING balance")
		log.Println(amount, balance_id, currency_id)
		tx.Rollback()
		return err
	}
	if newBalance < 0 {
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	return tx.Commit()
}

// Get exchange rates from server
func (r *WalletRepository) GetExchangeRates(ctx context.Context) (map[string]float64, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("server:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Create a new client
	client := pb.NewExchangeServiceClient(conn)

	// Create the context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call the GetExchangeRates method to retrieve the exchange rates from the server
	res, err := client.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get rates: %v", err)
	}
	// Extract the rates from the response
	rates := res.GetRates()

	return rates, nil
}

// GetExchangeRate retrieves the exchange rate between two currencies.
func (r *WalletRepository) GetExchangeRate(ctx context.Context, from string, to string) (float32, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("server:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Create a new client
	client := pb.NewExchangeServiceClient(conn)

	// Create the context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Make a new Currency request
	req := &pb.CurrencyRequest{
		FromCurrency: from,
		ToCurrency:   to,
	}

	// Call the GetExchangeRateForCurrency method to retrieve the exchange rate between two currencies from the server
	res, err := client.GetExchangeRateForCurrency(ctx, req)
	if err != nil {
		log.Println(from, to)
		log.Println("could not get rate: ", err)
		return 0, err
	}
	// Extract the rate from the response
	rate := res.GetRate()

	return rate, nil
}

// RegisterUser creates a new user account and wallets.
func (r *WalletRepository) RegisterUser(ctx context.Context, username string, email string, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// start a new transaction on the database
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// check if the username or email already exists
	var usernameScan string
	var emailScan string
	err = tx.QueryRow("SELECT username,email FROM mydb.users WHERE username = $1 OR email = $2", username, email).Scan(&usernameScan, &emailScan)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}

	if usernameScan == username {
		tx.Rollback()
		return errors.New("username already exists")
	}

	if emailScan == email {
		tx.Rollback()
		return errors.New("email already exists")
	}

	// insert the new user into the database
	hashedPassword := hashPassword(password)
	var lastInsertID int32
	err = r.db.QueryRowContext(ctx, "INSERT INTO mydb.users (username, email, password) VALUES ($1, $2, $3) RETURNING id", username, email, hashedPassword).Scan(&lastInsertID)
	if err != nil {
		return err
	}

	// create wallets for the new user
	var usWallet, eurWallet, rubWallet int32
	errUSD := r.db.QueryRowContext(ctx, "INSERT INTO mydb.wallets (user_id, name) VALUES ($1, $2) RETURNING id", lastInsertID, "USD WALLET").Scan(&usWallet)
	errRUB := r.db.QueryRowContext(ctx, "INSERT INTO mydb.wallets (user_id, name) VALUES ($1, $2) RETURNING id", lastInsertID, "RUB WALLET").Scan(&rubWallet)
	errEUR := r.db.QueryRowContext(ctx, "INSERT INTO mydb.wallets (user_id, name) VALUES ($1, $2) RETURNING id", lastInsertID, "EUR WALLET").Scan(&eurWallet)

	if errUSD != nil || errRUB != nil || errEUR != nil {
		tx.Rollback()
		return errors.New("failed to create wallets")
	}

	// create balances for the new wallets
	_, errUSD = r.db.ExecContext(ctx, "INSERT INTO mydb.balances (balance,wallet_id, currency_id) VALUES ($1, $2, $3) ", 0, usWallet, 2)
	_, errRUB = r.db.ExecContext(ctx, "INSERT INTO mydb.balances (balance,wallet_id, currency_id) VALUES ($1, $2, $3) RETURNING id", 0, rubWallet, 1)
	_, errEUR = r.db.ExecContext(ctx, "INSERT INTO mydb.balances (balance,wallet_id, currency_id) VALUES ($1, $2, $3) RETURNING id", 0, eurWallet, 3)
	if errUSD != nil || errRUB != nil || errEUR != nil {
		tx.Rollback()
		return errors.New("failed to create balances")

	}
	return tx.Commit()
}

// Login verifies the user credentials and generates a JWT token.
func (r *WalletRepository) Login(ctx context.Context, username, password string) (Token, error) {
	var uid int32
	var hashedPassword string

	// Query the database for the user's ID and hashed password
	err := r.db.QueryRowContext(ctx, "SELECT id,password FROM mydb.users WHERE username = $1", username).
		Scan(&uid, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return Token{}, errors.New("invalid username or password")
		}
		return Token{}, err
	}

	if hashPassword(password) != hashedPassword {
		return Token{}, errors.New("invalid username or password")
	}

	// Generate JWT token with user ID and username
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      float64(uid),
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return Token{}, err
	}
	return Token{Token: tokenString}, nil
}

// hashPassword hashes a password using SHA-256.
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
