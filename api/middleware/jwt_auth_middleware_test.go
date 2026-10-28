package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/middleware"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/internal/tokenutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	testSecret  = "secret-key-12345"
	wrongSecret = "wrong-key-67890"
)

func setupTestRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.JwtAuthMiddleware(secret))
	router.GET("/protected", func(c *gin.Context) {
		userID := c.GetString("x-user-id")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})
	return router
}

func TestJwtAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	router := setupTestRouter(testSecret)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var errResp domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Equal(t, "Not authorized", errResp.Message)
}

func TestJwtAuthMiddleware_MalformedAuthorizationHeader(t *testing.T) {
	router := setupTestRouter(testSecret)

	testCases := []struct {
		name   string
		header string
	}{
		{"Empty Bearer", "Bearer"},
		{"Single word", "TokenOnly"},
		{"Too many parts", "Bearer token extra-part"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tc.header)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)

			var errResp domain.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &errResp)
			assert.NoError(t, err)
			assert.Equal(t, "Not authorized", errResp.Message)
		})
	}
}

func TestJwtAuthMiddleware_InvalidToken(t *testing.T) {
	router := setupTestRouter(testSecret)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var errResp domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, errResp.Message)
}

func TestJwtAuthMiddleware_WrongSecretToken(t *testing.T) {
	router := setupTestRouter(testSecret)

	user := &domain.User{
		ID:    primitive.NewObjectID(),
		Name:  "Test User",
		Email: "test@example.com",
	}

	token, err := tokenutil.CreateAccessToken(user, wrongSecret, 2)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJwtAuthMiddleware_ValidToken(t *testing.T) {
	router := setupTestRouter(testSecret)

	userID := primitive.NewObjectID()
	user := &domain.User{
		ID:    userID,
		Name:  "Authorized User",
		Email: "auth@example.com",
	}

	token, err := tokenutil.CreateAccessToken(user, testSecret, 2)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, userID.Hex(), resp["user_id"])
}
