package controller

import (
	//"fmt"

	"net/http"
	"strconv"

	//	"sync/atomic"

	"simple-demo/repository"
	"simple-demo/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

//var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	db := repository.GETInstanceDB()
	uinfo, err := db.InsertUser(username, password)

	if err == gorm.ErrDuplicatedKey {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		token, err := util.NewToken(uinfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(uinfo.Id),
			Token:    token,
		})
	}

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	db := repository.GETInstanceDB()
	uinfo, err := db.LoginUser(username, password)

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else if uinfo.Id == -1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "wrong password"},
		})
	} else {
		token, err := util.NewToken(uinfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(uinfo.Id),
			Token:    token,
		})
	}
}

func UserInfo(c *gin.Context) {
	s_uid := c.Query("user_id")
	uid, err := strconv.ParseInt(s_uid, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "wrong uid"},
		})
		return
	}

	var uinfo repository.User
	db := repository.GETInstanceDB()
	uinfo, err = db.ShowUser(int(uid))

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "user not found"},
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		user := User{}
		user.Id = uid
		user.Name = uinfo.Username
		user.FollowCount = 0
		user.FollowerCount = 0
		user.IsFollow = false

		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	}

}
