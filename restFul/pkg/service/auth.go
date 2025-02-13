package service

import (
	"errors"
	database "service_api/pkg/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct{}

const (
	signingKey = "alskjjiopJ#OIjjkaslkj13dajlk159"
)

type authUserClaim struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := database.DataSource.GetUser(username, password)

	if err != nil {
		return "", err
	}

	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &authUserClaim{
		UserID: int(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(token string) (int, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &authUserClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid sign medthod")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := jwtToken.Claims.(*authUserClaim); !ok {
		return 0, errors.New("invalid claims")
	} else {
		return claims.UserID, nil
	}
}
