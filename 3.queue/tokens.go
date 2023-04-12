package main

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// TODO: token revoke after use

type _tokenServer struct {
	secret []byte
	// expireTimeInterval time.Time
}

type MyClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

func newTokenServer(salt string) *_tokenServer {
	return &_tokenServer{
		secret: []byte(salt),
		// expireTimeInterval: expireTime,
	}
}

func (ts *_tokenServer) generateToken(email string) (string, error) {
	claims := MyClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secret)
}

func (ts *_tokenServer) verifyToken(tokenString string) (string, bool) {
	token, _ := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return ts.secret, nil
	})

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims.Foo, true
	} else {
		return "", false
	}
}
