package controller

import (
	"DouYin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []service.VideoWithUser `json:"video_list"`
}

func Publish(c *gin.Context) {
	data, err := c.FormFile("data")
	if err != nil {
		log.Printf("get video data error:%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title := c.PostForm("title")
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoService := GetPublishVideo()
	err = videoService.Publish(data, userId, title)
	if err != nil {
		log.Printf("videoService.Publish() error：%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

func PublishList(c *gin.Context) {
	user_Id, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(user_Id, 10, 64)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoService := GetPublishVideo()
	list, err := videoService.GetVideosByAuthorid(userId, curId)
	if err != nil {
		log.Printf("videoService.GetVideosByAuthorid(%v)error：%v\n", userId, err)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "get video list error"},
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  Response{StatusCode: 0},
		VideoList: list,
	})
}

func GetPublishVideo() service.VideoWithUser {
	var user service.User
	var video service.VideoWithUser
	video.User = &user
	return video
}
