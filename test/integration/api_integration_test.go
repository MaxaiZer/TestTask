package integration

import (
	"bytes"
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"project/src/dto"
	"testing"
)

func Test_CreateTokens_WhenUserExists_ShouldReturn200(t *testing.T) {

	gin := setupRoutesForTests()
	req, _ := http.NewRequest("GET", "/auth/tokens?user_id=1", nil)
	w := httptest.NewRecorder()

	gin.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}

func Test_CreateTokens_WhenUserDoesntExists_ShouldReturn400(t *testing.T) {

	gin := setupRoutesForTests()
	req, _ := http.NewRequest("GET", "/auth/tokens?user_id=15", nil)
	w := httptest.NewRecorder()

	gin.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func Test_RefreshTokens_WhenInvalidTokens_ShouldReturn400(t *testing.T) {

	gin := setupRoutesForTests()

	tokenPair := dto.TokenPair{
		AccessToken:  "1",
		RefreshToken: "2",
	}
	body, _ := json.Marshal(tokenPair)

	req, _ := http.NewRequest("GET", "/auth/refresh", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	gin.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func Test_RefreshTokens_WhenValidTokens_ShouldReturn200(t *testing.T) {

	gin := setupRoutesForTests()

	req, _ := http.NewRequest("GET", "/auth/tokens?user_id=1", nil)
	w := httptest.NewRecorder()
	gin.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", "/auth/refresh", w.Body)
	gin.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}

func Test_RefreshTokens_WhenAlreadyUsedRefreshToken_ShouldReturn400(t *testing.T) {

	gin := setupRoutesForTests()

	req, _ := http.NewRequest("GET", "/auth/tokens?user_id=1", nil)
	w := httptest.NewRecorder()
	gin.ServeHTTP(w, req)

	bodyBytes := w.Body.Bytes()

	req, _ = http.NewRequest("GET", "/auth/refresh", bytes.NewBuffer(bodyBytes))
	w = httptest.NewRecorder()
	gin.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	req, _ = http.NewRequest("GET", "/auth/refresh", bytes.NewBuffer(bodyBytes))
	w = httptest.NewRecorder()
	gin.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}
