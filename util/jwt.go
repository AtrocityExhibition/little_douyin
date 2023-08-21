package util

import (
	"simple-demo/repository"
	"simple-demo/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.JwtSecret)

type MyClaims struct {
	Userid int `json:"user_id"`
	jwt.StandardClaims
}

func NewToken(user repository.User) (string, error) {
	c := MyClaims{
		user.Id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "AtrocityExhibition",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func Parsetoken(token string) int {
	tokenClaims, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*MyClaims)
		if ok && tokenClaims.Valid {
			return claims.Userid
		}
	}

	return -1
}
