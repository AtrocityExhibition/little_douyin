package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"simple-demo/util"
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

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	s_uid := c.Query("user_id")

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
		s_authorid := strconv.Itoa(v.Author_id)
		var tempvideoinfo Video
		tempvideoinfo.Id = int64(v.Id)
		tempvideoinfo.Author.Id = int64(v.Author_id)
		tempvideoinfo.Author.Name = s_authorid
		tempvideoinfo.Author.IsFollow = false
		tempvideoinfo.IsFavorite = false
		tempvideoinfo.Title = v.Title
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
