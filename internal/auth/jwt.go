// Package auth provides JWT token generation and validation for authentication.
package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// TokenExpiration defines the JWT token expiration time (8 hours)
	TokenExpiration = 8 * time.Hour
	
	// SecretKey is used to sign JWT tokens
	secretKey = "your-secret-key-change-in-production"
)

// Claims represents the JWT claims structure
type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// HashPassword creates a SHA-256 hash of the password
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// GenerateToken creates a new JWT token for the given password
func GenerateToken(password string) (string, error) {
	passwordHash := HashPassword(password)
	
	claims := Claims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateToken validates the JWT token and checks if the password hash matches
func ValidateToken(tokenString string, currentPassword string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("token is empty")
	}
	
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})
	
	if err != nil {
		return false, err
	}
	
	if !token.Valid {
		return false, errors.New("invalid token")
	}
	
	// Check if the password has changed by comparing hashes
	currentHash := HashPassword(currentPassword)
	if claims.PasswordHash != currentHash {
		return false, errors.New("password has changed")
	}
	
	return true, nil
}

