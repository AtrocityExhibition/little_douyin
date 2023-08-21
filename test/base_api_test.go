package test

import (
	"fmt"
	"math/rand"
	"net/http"
	"simple-demo/repository"
	"simple-demo/util"
	"testing"
	"time"
)

func TestMainPage(t *testing.T) {
	e := newExpect(t)
	e.GET("/douyin/").Expect().Status(http.StatusOK).JSON().Object()
}

func TestFeed(t *testing.T) {
	e := newExpect(t)

	feedResp := e.GET("/douyin/feed/").WithQuery("latest_time", "1692618767539").Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().Equal(0)
	feedResp.Value("video_list").Array().Length().Gt(0)

	for _, element := range feedResp.Value("video_list").Array().Iter() {
		video := element.Object()
		video.ContainsKey("id")
		video.ContainsKey("author")
		video.Value("play_url").String().NotEmpty()
		video.Value("cover_url").String().NotEmpty()
		//1
	}
}

func TestUserAction(t *testing.T) {
	e := newExpect(t)

	rand.Seed(time.Now().UnixNano())
	registerValue := fmt.Sprintf("douyin%d", rand.Intn(65536))

	registerResp := e.POST("/douyin/user/register/").
		WithQuery("username", registerValue).WithQuery("password", registerValue).
		WithFormField("username", registerValue).WithFormField("password", registerValue).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	registerResp.Value("status_code").Number().Equal(0)
	registerResp.Value("user_id").Number().Gt(0)
	registerResp.Value("token").String().Length().Gt(0)

	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", registerValue).WithQuery("password", registerValue).
		WithFormField("username", registerValue).WithFormField("password", registerValue).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	loginResp.Value("status_code").Number().Equal(0)
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().Length().Gt(0)

	token, _ := util.NewToken(repository.User{Id: 4})
	uid := 4
	userResp := e.GET("/douyin/user/").
		WithQuery("token", token).
		WithQuery("user_id", uid).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	userResp.Value("status_code").Number().Equal(0)
	userInfo := userResp.Value("user").Object()
	userInfo.NotEmpty()
	userInfo.Value("id").Number().Gt(0)
	userInfo.Value("name").String().Length().Gt(0)
	//1  r3f
}

func TestPublish(t *testing.T) {
	e := newExpect(t)
	token, _ := util.NewToken(repository.User{Id: 4})

	publishResp := e.POST("/douyin/publish/action/").
		WithMultipart().
		WithFile("data", "../public/bear.mp4").
		WithFormField("token", token).
		WithFormField("title", "Bear").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	publishResp.Value("status_code").Number().Equal(0)

	publishListResp := e.GET("/douyin/publish/list/").
		WithQuery("user_id", 4).WithQuery("token", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	publishListResp.Value("status_code").Number().Equal(0)
	publishListResp.Value("video_list").Array().Length().Gt(0)

	for _, element := range publishListResp.Value("video_list").Array().Iter() {
		video := element.Object()
		video.ContainsKey("id")
		video.ContainsKey("author")
		video.Value("play_url").String().NotEmpty()
		video.Value("cover_url").String().NotEmpty()
	}
}
