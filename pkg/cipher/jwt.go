package cipher

import (
	"crypto/ecdsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAuth struct {
	key *ecdsa.PrivateKey
}

func NewJwtAuth(key *ecdsa.PrivateKey) *JwtAuth {
	return &JwtAuth{key: key}
}

func (ja *JwtAuth) NewJwtTOken(expireTime time.Time) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "kube-eip",
		ExpiresAt: jwt.NewNumericDate(expireTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenStr, err := token.SignedString(ja.key)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (ja *JwtAuth) ValidateJwtToken(tokenStr string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return ja.key.Public(), nil
	})

	switch {
	case token.Valid:
		return true, nil
	default:
		return false, err
	}
}
