package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"test-task/internal/dto"
	"testing"
	"time"
)

func TestGetBalance_WhenInvalidUUID_ShouldReturn400(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/wallets/123", nil)
	w := httptest.NewRecorder()

	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBalance_WhenWalletDoesntExist_ShouldReturn404(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/wallets/123e4567-e89b-12d3-a456-426614174000", nil)
	w := httptest.NewRecorder()

	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetBalance_WhenWalletExists_ShouldReturn200(t *testing.T) {

	walletId := "11111111-1111-1111-1111-111111111111"
	balance, err := getBalance(ginEngine, walletId)

	assert.NoError(t, err)
	assert.Equal(t, 555.5, balance)
}

func TestOperation_WhenInvalidOperation_ShouldReturn400(t *testing.T) {

	op := dto.WalletOperation{
		WalledID:      "11111111-1111-1111-1111-111111111111",
		OperationType: "someRandomOperation",
		Amount:        15,
	}
	body, _ := json.Marshal(op)

	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOperation_WhenWalletDoesntExist_ShouldReturn404(t *testing.T) {

	op := dto.WalletOperation{
		WalledID:      "123e4567-e89b-12d3-a456-426614174000",
		OperationType: "DEPOSIT",
		Amount:        15,
	}
	body, _ := json.Marshal(op)

	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOperation_WhenBalanceBecomeLessThanZero_ShouldReturn400(t *testing.T) {

	op := dto.WalletOperation{
		WalledID:      "11111111-1111-1111-1111-111111111111",
		OperationType: "WITHDRAW",
		Amount:        9999999999,
	}
	body, _ := json.Marshal(op)

	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWithdraw(t *testing.T) {

	engine := setupRoutesForTests()

	walletId := "11111111-1111-1111-1111-111111111111"
	prevBalance, err := getBalance(engine, walletId)
	assert.NoError(t, err)

	op := dto.WalletOperation{
		WalledID:      walletId,
		OperationType: "WITHDRAW",
		Amount:        10,
	}
	err = runOperation(engine, op)
	assert.NoError(t, err)

	balance, err := getBalance(engine, op.WalledID)
	assert.NoError(t, err)
	assert.Equal(t, prevBalance-op.Amount, balance)
}

func TestDeposit(t *testing.T) {

	engine := setupRoutesForTests()

	walletId := "11111111-1111-1111-1111-111111111111"
	prevBalance, err := getBalance(engine, walletId)
	assert.NoError(t, err)

	op := dto.WalletOperation{
		WalledID:      walletId,
		OperationType: "DEPOSIT",
		Amount:        10,
	}
	err = runOperation(engine, op)
	assert.NoError(t, err)

	balance, err := getBalance(engine, op.WalledID)
	assert.NoError(t, err)
	assert.Equal(t, prevBalance+op.Amount, balance)
}

func TestConcurrentWalletOperations(t *testing.T) {

	walletID := "11111111-1111-1111-1111-111111111111"

	duration := time.Second
	numRequests := 700
	depositAmount := 1000.0

	prevBalance, err := getBalance(ginEngine, walletID)
	assert.NoError(t, err)

	start := time.Now()
	log.Infof("Running %d requests...", numRequests)
	limiter := rate.NewLimiter(rate.Every(duration/time.Duration(numRequests)), 1)

	wg := sync.WaitGroup{}
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			op := dto.WalletOperation{
				WalledID:      walletID,
				OperationType: "DEPOSIT",
				Amount:        depositAmount,
			}
			err := limiter.Wait(context.Background())
			assert.NoError(t, err)

			err = runOperation(ginEngine, op)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()
	log.Infof("Handled %d requests in %fs, avg per request: %fs", numRequests,
		time.Since(start).Seconds(),
		time.Since(start).Seconds()/float64(numRequests))

	balance, err := getBalance(ginEngine, walletID)
	assert.NoError(t, err)
	assert.Equal(t, prevBalance+float64(numRequests)*depositAmount, balance)
}

func TestConcurrentWalletOperations_WithExceedingRateLimit(t *testing.T) {

	walletID := "11111111-1111-1111-1111-111111111111"

	duration := time.Second
	numRequests := 2000
	has429 := false

	start := time.Now()
	log.Infof("Running %d requests...", numRequests)
	limiter := rate.NewLimiter(rate.Every(duration/time.Duration(numRequests)), 1)

	wg := sync.WaitGroup{}
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			err := limiter.Wait(context.Background())
			assert.NoError(t, err)

			r := rand.Int() % 3
			if r == 0 {
				op := dto.WalletOperation{WalledID: walletID, OperationType: "DEPOSIT", Amount: 0.01}
				err = runOperation(ginEngine, op)
			} else if r == 1 {
				op := dto.WalletOperation{WalledID: walletID, OperationType: "WITHDRAW", Amount: 0.01}
				err = runOperation(ginEngine, op)
			} else if r == 2 {
				_, err = getBalance(ginEngine, walletID)
			}

			if err != nil {
				if strings.Contains(err.Error(), "429") {
					has429 = true
				} else {
					assert.Fail(t, "api returned not 429 error")
				}
			}
		}(i)
	}
	wg.Wait()
	log.Infof("Handled %d requests in %fs, avg per request: %fs", numRequests,
		time.Since(start).Seconds(),
		time.Since(start).Seconds()/float64(numRequests))

	assert.True(t, has429)
}

func getBalance(gin *gin.Engine, walletID string) (float64, error) {
	req, _ := http.NewRequest("GET", "/api/v1/wallets/"+walletID, nil)
	w := httptest.NewRecorder()
	gin.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", w.Code)
	}

	var response dto.WalletBalance
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return 0, err
	}

	return response.Balance, nil
}

func runOperation(engine *gin.Engine, op dto.WalletOperation) error {
	body, _ := json.Marshal(op)
	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", w.Code)
	}

	return nil
}
