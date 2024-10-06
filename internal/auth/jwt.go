package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret   string
	issuer   string
	audience string
}

func NewJWTAuthenticator(secret, issuer, audience string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:   secret,
		issuer:   issuer,
		audience: audience,
	}
}

func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(a.issuer),
		jwt.WithAudience(a.audience),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
