package auth

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type JWTService struct {
	Auth *jwtauth.JWTAuth
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		Auth: jwtauth.New("HS256", []byte(secret), nil),
	}
}

func (j *JWTService) CreateToken(userID int64) (string, error) {
	claims := map[string]interface{}{
		"user_id": userID,
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
	}

	_, tokenString, err := j.Auth.Encode(claims)
	return tokenString, err
}
