package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenManager struct {
	signingKey string
}

type TokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

func NewTokenManager(signingKey string) (*TokenManager, error) {
	if signingKey == "" {
		return nil, fmt.Errorf("empty signing key")
	}

	return &TokenManager{signingKey: signingKey}, nil
}

func (m *TokenManager) GenerateAccessToken(userID int) (string, error) {
	// Срок действия access token - 15 минут
	claims := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.signingKey))
}

func (m *TokenManager) GenerateRefreshToken() (string, error) {
	// Генерируем случайный refresh token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(720 * time.Hour).Unix() // 30 дней
	claims["iat"] = time.Now().Unix()

	return token.SignedString([]byte(m.signingKey))
}

func (m *TokenManager) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, fmt.Errorf("error getting user claims from token")
	}

	return claims.UserID, nil
}
