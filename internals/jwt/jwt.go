package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService defines the interface for token operations.
type TokenService interface {
	GenerateToken(username string) (string, error)
	Validate(tokenString string) (bool, error)
}

type JWT struct {
	Secret []byte
}

// NewJWT creates a new JWT instance with the provided secret.
func NewJWT(secret string) *JWT {
	return &JWT{Secret: []byte(secret)}
}

func (s *JWT) GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(s.Secret)
}

func (s *JWT) Validate(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.Secret, nil
	})
	if err != nil {
		return false, err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}
	return false, nil
}