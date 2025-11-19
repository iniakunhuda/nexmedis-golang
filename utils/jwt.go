package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	ClientID uuid.UUID `json:"client_id"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// InitJWT initializes the JWT secret
func InitJWT() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-in-production"
	}
	jwtSecret = []byte(secret)
}

// GenerateJWT generates a JWT token for a client
func GenerateJWT(clientID uuid.UUID, email string) (string, error) {
	if len(jwtSecret) == 0 {
		InitJWT()
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	
	// Get custom expiration from env if set
	if expStr := os.Getenv("JWT_EXPIRATION"); expStr != "" {
		if exp, err := time.ParseDuration(expStr); err == nil {
			expirationTime = time.Now().Add(exp)
		}
	}

	claims := &JWTClaims{
		ClientID: clientID,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	if len(jwtSecret) == 0 {
		InitJWT()
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshJWT generates a new JWT token from an existing valid token
func RefreshJWT(tokenString string) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}

	return GenerateJWT(claims.ClientID, claims.Email)
}
