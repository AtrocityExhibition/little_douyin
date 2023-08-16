package controller

import (
	//"fmt"
	"fmt"
	"net/http"
	"strconv"

	//	"sync/atomic"

	"simple-demo/repository"

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

	var uid int
	db := repository.GETInstanceDB()
	err := db.InsertUser(username, password, &uid)

	if err == gorm.ErrDuplicatedKey {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(uid),
			Token:    username + password,
		})
	}

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var uid int
	db := repository.GETInstanceDB()
	err := db.LoginUser(username, password, &uid)
	fmt.Println(uid)

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else if uid == -1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "wrong password"},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(uid),
			Token:    username + password,
		})
	}
}

func UserInfo(c *gin.Context) {
	s_uid := c.Query("user_id")
	token := c.Query("token")
	uid, err := strconv.ParseInt(s_uid, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "wrong uid"},
		})
		return
	}

	var uinfo repository.User
	db := repository.GETInstanceDB()
	err = db.ShowUser(int(uid), &uinfo)

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "user not found"},
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else if uinfo.Username+uinfo.Password != token {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "authenticate failed"},
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
