package controller

import (
	"DouYin/repository"
	"DouYin/service"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []service.VideoWithUser `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	uid := c.GetInt("userId")

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	playUrl := "https://douyin-duu.oss-cn-beijing.aliyuncs.com/" + finalName
	CoverUrl := playUrl + "?x-oss-process=video/snapshot,t_0,f_jpg,w_800,h_600"

	// 上传到oss
	file, err := data.Open()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return
	}
	err = util.Bucket.PutObject(finalName, bytes.NewReader(buf.Bytes()))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//saveFile := filepath.Join("./public/", finalName)
	//if err := c.SaveUploadedFile(data, saveFile); err != nil {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}

	var tempvideo repository.Video
	title := c.PostForm("title")
	tempvideo.Title = title
	tempvideo.Author_id = uid
	tempvideo.PlayUrl = playUrl
	tempvideo.CoverUrl = CoverUrl

	db := repository.GETInstanceDB()
	err = db.InsertVideo(tempvideo)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
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
