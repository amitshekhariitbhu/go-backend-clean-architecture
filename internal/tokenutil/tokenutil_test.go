package tokenutil_test

import (
	"testing"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/internal/tokenutil"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	testSecret      = "test-secret-key-12345"
	wrongSecret     = "wrong-secret-key-67890"
	tokenExpiryHour = 2
)

func TestCreateAccessToken_Success(t *testing.T) {
	userID := primitive.NewObjectID()
	user := &domain.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	token, err := tokenutil.CreateAccessToken(user, testSecret, tokenExpiryHour)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	authorized, err := tokenutil.IsAuthorized(token, testSecret)
	assert.NoError(t, err)
	assert.True(t, authorized)

	extractedID, err := tokenutil.ExtractIDFromToken(token, testSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID.Hex(), extractedID)
}

func TestCreateAccessToken_WrongSecret(t *testing.T) {
	userID := primitive.NewObjectID()
	user := &domain.User{
		ID:    userID,
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	token, err := tokenutil.CreateAccessToken(user, testSecret, tokenExpiryHour)
	assert.NoError(t, err)

	authorized, err := tokenutil.IsAuthorized(token, wrongSecret)
	assert.Error(t, err)
	assert.False(t, authorized)

	_, err = tokenutil.ExtractIDFromToken(token, wrongSecret)
	assert.Error(t, err)
}

func TestCreateRefreshToken_Success(t *testing.T) {
	userID := primitive.NewObjectID()
	user := &domain.User{
		ID:    userID,
		Name:  "Refresh User",
		Email: "refresh@example.com",
	}

	refreshToken, err := tokenutil.CreateRefreshToken(user, testSecret, tokenExpiryHour)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	authorized, err := tokenutil.IsAuthorized(refreshToken, testSecret)
	assert.NoError(t, err)
	assert.True(t, authorized)

	extractedID, err := tokenutil.ExtractIDFromToken(refreshToken, testSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID.Hex(), extractedID)
}

func TestIsAuthorized_InvalidToken(t *testing.T) {
	authorized, err := tokenutil.IsAuthorized("invalid.token.string", testSecret)
	assert.Error(t, err)
	assert.False(t, authorized)
}

func TestExtractIDFromToken_InvalidToken(t *testing.T) {
	extractedID, err := tokenutil.ExtractIDFromToken("invalid.token.string", testSecret)
	assert.Error(t, err)
	assert.Empty(t, extractedID)
}

func TestToken_Expired(t *testing.T) {
	userID := primitive.NewObjectID()
	user := &domain.User{
		ID:    userID,
		Name:  "Expired User",
		Email: "expired@example.com",
	}

	// Token created with negative expiry (already expired)
	token, err := tokenutil.CreateAccessToken(user, testSecret, -1)
	assert.NoError(t, err)

	authorized, err := tokenutil.IsAuthorized(token, testSecret)
	assert.Error(t, err)
	assert.False(t, authorized)
	assert.Contains(t, err.Error(), "Token is expired")
}
