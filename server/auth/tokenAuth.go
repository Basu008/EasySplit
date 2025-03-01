package auth

import (
	"encoding/base64"
	"os"

	"github.com/Basu008/EasySplit.git/server/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenAuth interface {
	SignToken() (string, error)
	VerifyToken(string) (*UserClaim, error)
}

type TokenAuthentication struct {
	Config    *config.TokenAuthConfig
	UserClaim *UserClaim
}

func (t *TokenAuthentication) SignToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t.UserClaim)
	secretKey := []byte(t.Config.JWTSignKey)
	tokenString, _ := token.SignedString(secretKey)
	return base64.StdEncoding.EncodeToString([]byte(tokenString)), nil
}

func (t *TokenAuthentication) VerifyToken(tokenString string) (*UserClaim, error) {
	decodedString, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, err
	}
	var claim UserClaim
	token, err := jwt.ParseWithClaims(string(decodedString), &claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return &claim, err
	}
	if !token.Valid {
		return &claim, err
	}
	return &claim, err
}

func NewTokenAuthentication(c *config.TokenAuthConfig) *TokenAuthentication {
	return &TokenAuthentication{Config: c}
}

type UserClaim struct {
	ID          uint   `json:"id"`
	Type        string `json:"type"`
	PhoneNumber string `json:"phone_no"`
	jwt.RegisteredClaims
}
