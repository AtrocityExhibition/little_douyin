package repository

import (
	"DouYin/config"
	"DouYin/util"
	"io"
	"io/ioutil"
	"log"
	"path"
	"time"
)

type Video struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	PlayUrl     string `json:"play_url"`
	CoverUrl    string `json:"cover_url"`
	PublishTime time.Time
	Title       string `json:"title"`
}

func QueryVideosByAuthorId(authorId int64) ([]Video, error) {
	var data []Video
	result := Db.Where(&Video{AuthorId: authorId}).Find(&data)
	// return Video struct
	return data, result.Error
}

func QueryVideoByVideoId(videoId int64) (Video, error) {
	var video Video
	video.Id = videoId
	result := Db.First(&video)
	return video, result.Error

}

func QueryVideosByLastTime(lastTime time.Time) ([]Video, error) {
	videos := make([]Video, config.VideoCount)
	result := Db.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(config.VideoCount).Find(&videos)
	return videos, result.Error
}

// VideoSFTP upload video to sftp
func VideoSFTP(file io.Reader, videoName string) error {
	// upload video to sftp
	// create remote file
	ftpFile, err := util.SftpClient.Create(path.Join("home/zhouyx/video/", videoName+".mp4")) // 这里的remotePath是sftp根目录下的目录，是目录不是文件名
	if nil != err {
		log.Println("sftpClient.Create error", err)
		return err
	}
	defer ftpFile.Close()
	// wirte remote file
	fileByte, err := ioutil.ReadAll(file)
	if nil != err {
		log.Println("ioutil.ReadAll error", err)
		return err
	}
	ftpFile.Write(fileByte)
	return nil
}

func SaveVideoToDB(videoName string, imageName string, authorId int64, title string) error {
	var video Video
	video.PublishTime = time.Now()
	video.PlayUrl = config.PlayUrlPrefix + videoName + ".mp4"
	video.CoverUrl = config.CoverUrlPrefix + imageName + ".jpg"
	video.AuthorId = authorId
	video.Title = title
	result := Db.Create(&video)
	return result.Error
}
