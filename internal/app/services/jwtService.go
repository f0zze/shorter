package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenExp  = time.Hour * 3
	SecretKey = "authKey"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokeString, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return "", err
	}

	return tokeString, nil
}

func GetUserID(tokenString string) (string, bool) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return "", false
	}

	fmt.Println("Token maybe have an error ", err.Error())
	if !token.Valid {
		fmt.Println("Token is not valid")
		return "", false
	}

	fmt.Println("Token is valid")
	return claims.UserID, true
}
