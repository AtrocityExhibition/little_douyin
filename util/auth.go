package util

import (
	"DouYin/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// Auth if token is right, put userId in context, and return ok, flag is true for Get
func Auth(flag bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth string
		if flag == true {
			auth = c.Query("token")
		} else {
			auth = c.Request.PostFormValue("token")
		}
		if len(auth) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		auth = strings.Fields(auth)[1]
		token, err := parseToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
		}
		c.Set("userId", token.Id)
		c.Next()
	}
}

// AuthWithoutLogin if bring token, put it in userId, els put 0 in it
func AuthWithoutLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Query("token")
		var userId string
		if len(auth) == 0 {
			userId = "0"
		} else {
			auth = strings.Fields(auth)[1]
			token, err := parseToken(auth)
			if err != nil {
				c.Abort()
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			} else {
				userId = token.Id
			}
		}
		c.Set("userId", userId)
		c.Next()
	}
}

// parseToken
func parseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return config.Secret, nil
	})
	if err == nil && jwtToken != nil {
		jwtclaim, flag := jwtToken.Claims.(*jwt.StandardClaims)
		if flag && jwtToken.Valid {
			return jwtclaim, nil
		}
	}
	return nil, err
}
