package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
	"sync"
	"test-task/internal/entities"
	"test-task/internal/errors"
	"time"
)

type operationName string

const (
	withdraw operationName = "WITHDRAW"
	deposit  operationName = "DEPOSIT"
)

type WalletOperation struct {
	walletID string
	name     operationName
	amount   float64
}

func NewWalletOperation(walletID string, operation string, amount float64) (*WalletOperation, error) {

	if walletID == "" {
		return nil, fmt.Errorf("walletID is empty")
	}

	if operation == "" {
		return nil, fmt.Errorf("operation is empty")
	}

	if _, err := uuid.Parse(walletID); err != nil {
		return nil, fmt.Errorf("walletID is not uuid")
	}

	if operationName(operation) != withdraw && operationName(operation) != deposit {
		return nil, fmt.Errorf("invalid operation name, expected: %s or %s", withdraw, deposit)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("operation amount must be greater than zero")
	}

	return &WalletOperation{walletID, operationName(operation), amount}, nil
}

type walletsRepository interface {
	GetById(ctx context.Context, id string) (entities.Wallet, error)
	ChangeBalance(ctx context.Context, id string, delta float64) error
}

type walletLimiter struct {
	limiter         *rate.Limiter
	lastRequestTime time.Time
}

type WalletsService struct {
	wallets       walletsRepository
	limiters      map[string]*walletLimiter
	mu            sync.Mutex
	cancelCleanup context.CancelFunc
}

func NewWalletsService(wallets walletsRepository) *WalletsService {
	service := &WalletsService{wallets: wallets, limiters: make(map[string]*walletLimiter)}

	ctx, cancel := context.WithCancel(context.Background())
	go service.limitersCleanup(ctx)
	service.cancelCleanup = cancel
	return service
}

func (s *WalletsService) GetBalance(ctx context.Context, id string) (float64, error) {

	if !s.allowWalletOperation(id) {
		return 0, errors.TooManyRequests
	}

	wallet, err := s.wallets.GetById(ctx, id)
	return wallet.Balance, err
}

func (s *WalletsService) RunOperation(ctx context.Context, operation WalletOperation) error {

	if !s.allowWalletOperation(operation.walletID) {
		return errors.TooManyRequests
	}

	switch operation.name {
	case withdraw:
		return s.wallets.ChangeBalance(ctx, operation.walletID, -operation.amount)
	case deposit:
		return s.wallets.ChangeBalance(ctx, operation.walletID, operation.amount)
	default:
		return fmt.Errorf("%w: %s", errors.UnsupportedOperation, operation.name)
	}
}

func (s *WalletsService) Close() {
	s.cancelCleanup()
}

func (s *WalletsService) allowWalletOperation(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, ok := s.limiters[id]
	if !ok {
		limiter = &walletLimiter{
			limiter:         rate.NewLimiter(1000, 5),
			lastRequestTime: time.Now(),
		}
		s.limiters[id] = limiter
	}
	return limiter.limiter.Allow()
}

func (s *WalletsService) limitersCleanup(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Minute):
			s.mu.Lock()
			for id, limiter := range s.limiters {
				if time.Since(limiter.lastRequestTime) > 5*time.Minute {
					delete(s.limiters, id)
				}
			}
			s.mu.Unlock()
		}
	}
}
