package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"test-task/internal/entities"
	errs "test-task/internal/errors"
)

const balanceNotNegativeCheck = "check_balance_non_negative"

type Wallets struct {
	db *sqlx.DB
}

func NewWalletsRepository(db *sqlx.DB) *Wallets {
	return &Wallets{db: db}
}

func (repo *Wallets) GetById(ctx context.Context, id string) (entities.Wallet, error) {
	var wallet entities.Wallet
	err := repo.db.GetContext(ctx, &wallet, "SELECT * FROM wallets WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wallet, fmt.Errorf("%w: wallet by id %s", errs.NotFound, id)
		}
		return wallet, err
	}

	return wallet, nil
}

func (repo *Wallets) ChangeBalance(ctx context.Context, id string, delta float64) error {

	res, err := repo.db.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 where id = $2", delta, id)
	if err != nil {
		if strings.Contains(err.Error(), balanceNotNegativeCheck) {
			return errs.InsufficientBalance
		}
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: wallet by id %s", errs.NotFound, id)
	}

	return nil
}
