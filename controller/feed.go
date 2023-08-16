package controller

import (
	"net/http"
	"strconv"
	"time"

	"simple-demo/repository"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var searchtime time.Time
	var nexttime time.Time
	var err error
	lasttime := c.Query("latest_time")
	if len(lasttime) != 0 {
		searchtime, err = time.Parse("2006-01-02 15:04:05", lasttime)
	} else {
		searchtime = time.Now()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "parse time error" + err.Error()},
		})
	}

	var r_videos []repository.Video
	db := repository.GETInstanceDB()
	r_videos, err = db.SearchVideosbyTime(searchtime)

	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else if len(r_videos) == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "video don't exists"},
		})
	} else {
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

		if len(r_videos) < 30 {
			nexttime = time.Now()
		} else {
			nexttime = r_videos[29].Createtime
		}

		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videos,
			NextTime:  nexttime.Unix(),
		})
	}
}
