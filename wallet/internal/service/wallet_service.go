package service

import (
	"context"
	"errors"
	"wallet/internal/repository"
)

type WalletServiceInterface interface {
	Deposit(ctx context.Context, uid int32, amount int32, currency string) error
	Withdraw(ctx context.Context, uid int32, amount int32, currency string) error
	GetBalance(ctx context.Context, username string) (map[string]int32, error)
	GetExchangeRates(ctx context.Context) (map[string]float64, error)
	GetExchangeRate(ctx context.Context, from string, to string) (float32, error)
	RegisterUser(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, username, password string) (repository.Token, error)
}

type WalletService struct {
	repo repository.WalletRepositoryInterface
}

func NewWalletService(repo repository.WalletRepositoryInterface) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) Deposit(ctx context.Context, uid int32, amount int32, currency string) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return s.repo.UpdateBalance(ctx, uid, amount, currency)
}

func (s *WalletService) Withdraw(ctx context.Context, uid int32, amount int32, currency string) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	return s.repo.UpdateBalance(ctx, uid, -amount, currency)
}

func (s *WalletService) GetBalance(ctx context.Context, username string) (map[string]int32, error) {
	return s.repo.GetBalance(ctx, username)
}

func (s *WalletService) GetExchangeRates(ctx context.Context) (map[string]float64, error) {
	return s.repo.GetExchangeRates(ctx)
}

func (s *WalletService) GetExchangeRate(ctx context.Context, from string, to string) (float32, error) {
	return s.repo.GetExchangeRate(ctx, from, to)
}

func (s *WalletService) RegisterUser(ctx context.Context, username string, email string, password string) error {
	return s.repo.RegisterUser(ctx, username, email, password)
}

func (s *WalletService) Login(ctx context.Context, username string, password string) (repository.Token, error) {
	return s.repo.Login(ctx, username, password)
}

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrWalletNotFound    = errors.New("wallet not found")
)
