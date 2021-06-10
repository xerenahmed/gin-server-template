package utils

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"time"
)

var (
	AccessSecret = []byte(os.Getenv("ACCESS_SECRET"))
)

func GenerateAccessToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24 * 4).Unix()

	tokenString, err := token.SignedString(AccessSecret)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}

	return tokenString, nil
}

func ParseToken(secret []byte, tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}
