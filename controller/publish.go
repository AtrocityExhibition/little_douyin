package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"simple-demo/repository"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	//todo: 鉴权信息保留在内存中, 是可以的, 但是不能用map, 好像可以用 sync.map
	/*
		if _, exist := usersLoginInfo[token]; !exist {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
			return
		}
	*/

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
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//todo: 改一下表结构, 必须有filename以生成play_url
	//todo: 考虑一下在token中间加入 id + usernmae
	title := c.PostForm("title")
	var tempvideo repository.Video
	tempvideo.Title = title
	tempvideo.Author_id = 12

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

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	s_uid := c.Query("user_id")
	//todo: token
	// token := c.Query("token")

	db := repository.GETInstanceDB()
	uid, err := strconv.ParseInt(s_uid, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	r_videos, err := db.SearchPublished(int(uid))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videos := []Video{}
	for _, v := range r_videos {
		//todo: author_name 直接存在视频表里还是在此处搜索用户信息, 这取决于是否实现拓展功能
		//todo: play_url 和 cover_url 修改
		//todo: 回应结构体数据待完善, 没有title
		s_authorid := strconv.Itoa(v.Author_id)
		var tempvideoinfo Video
		tempvideoinfo.Id = int64(v.Id)
		tempvideoinfo.Author.Id = int64(v.Author_id)
		tempvideoinfo.Author.Name = s_authorid
		tempvideoinfo.Author.IsFollow = false
		tempvideoinfo.IsFavorite = false
		tempvideoinfo.CoverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
		tempvideoinfo.PlayUrl = "https://www.w3schools.com/html/movie.mp4"
		videos = append(videos, tempvideoinfo)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
