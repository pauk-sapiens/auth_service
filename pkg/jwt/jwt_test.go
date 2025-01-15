package jwt

import (
	"auth/pkg/core/models"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Mock data
var testUser = models.User{
	ID:     1,
	Email:  "test@example.com",
	PWHash: []byte("blabla"),
}

var testApp = models.App{
	ID:     1,
	Secret: "supersecretkey",
}

func TestNewToken_Success(t *testing.T) {
	duration := time.Hour
	tokenString, err := NewToken(testUser, testApp, duration)

	assert.NoError(t, err, "expected no error while generating token")
	assert.NotEmpty(t, tokenString, "expected non-empty token string")

	// Parse token to validate claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(testApp.Secret), nil
	})

	assert.NoError(t, err, "expected no error while parsing token")
	assert.True(t, token.Valid, "expected valid token")

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok, "expected claims to be of type MapClaims")
	assert.Equal(t, float64(testUser.ID), claims["uid"], "expected user ID to match")
	assert.Equal(t, testUser.Email, claims["email"], "expected email to match")
	assert.Equal(t, float64(testApp.ID), claims["appID"], "expected app ID to match")
}

func TestNewToken_InvalidSecret(t *testing.T) {
	// Test with an empty secret to simulate an error
	invalidApp := models.App{
		ID:     2,
		Secret: "",
	}
	duration := time.Hour

	tokenString, err := NewToken(testUser, invalidApp, duration)

	assert.Error(t, err, "expected invalid key error")
	assert.Empty(t, tokenString, "", "expected empty token string")
}
