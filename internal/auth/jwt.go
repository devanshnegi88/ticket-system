package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken is returned for any token that fails parsing,
// signature verification, or has expired.
var ErrInvalidToken = errors.New("invalid or expired token")

// Claims is the JWT payload. UserID is the source of truth used for all
// ownership checks throughout the application.
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager creates and verifies HS256-signed JWTs.
type JWTManager struct {
	secret []byte
	expiry time.Duration
}

// NewJWTManager builds a JWTManager from a secret string and token
// lifetime.
func NewJWTManager(secret string, expiry time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), expiry: expiry}
}

// Generate issues a new signed JWT for the given user.
func (m *JWTManager) Generate(userID uint, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Verify parses and validates a token string, returning its claims if
// valid, or ErrInvalidToken otherwise.
func (m *JWTManager) Verify(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
