package middleware

import (
	"net/http"
	"simple-demo/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type failresponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 参数为 token + id, Get方法, 检验token是否正确
func AuthWithId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		s_userid := ctx.Query("user_id")
		userid, err := strconv.ParseInt(s_userid, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, failresponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
		}

		uid := util.Parsetoken(token)
		if uid == -1 || uid != int(userid) {
			ctx.AbortWithStatusJSON(http.StatusOK, failresponse{
				StatusCode: 1,
				StatusMsg:  "wrong token",
			})
		}
	}
}

// 参数为token, Get方法, 检验token并将id放入上下文中
func AuthGetId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")

		uid := util.Parsetoken(token)
		if uid == -1 {
			ctx.AbortWithStatusJSON(http.StatusOK, failresponse{
				StatusCode: 1,
				StatusMsg:  "wrong token",
			})
		}

		ctx.Set("userId", uid)
	}
}

// 参数为token, Post方法, 检验token并将id放入上下文中
func AuthPostId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.PostForm("token")

		uid := util.Parsetoken(token)
		if uid == -1 {
			ctx.AbortWithStatusJSON(http.StatusOK, failresponse{
				StatusCode: 1,
				StatusMsg:  "wrong token",
			})
		}

		ctx.Set("userId", uid)
	}
}
