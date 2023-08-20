package controller

import (
	"DouYin/config"
	"DouYin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []service.VideoWithUser `json:"video_list,omitempty"`
	NextTime  int64                   `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	var lastTime time.Time
	var err error
	inputTime := c.Query("latest_time")
	if len(inputTime) != 0 {
		lastTime, err = time.Parse(config.DateTime, inputTime)
	} else {
		lastTime = time.Now()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "parse time error" + err.Error()},
		})
	}
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoService := GetFeedVideo()
	feed, nextTime, err := videoService.GetVideosByLastTime(lastTime, userId)
	if err != nil {
		log.Printf("videoService.Feed() error：%v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "get video data error"},
		})
		return
	} else if len(feed) == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "video don't exists"},
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  nextTime.Unix(),
	})
}

func GetFeedVideo() service.VideoWithUser {
	var user service.User
	var video service.VideoWithUser
	video.User = &user
	return video
}
