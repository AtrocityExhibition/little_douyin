package service

import (
	"DouYin/config"
	"DouYin/repository"
	"DouYin/util"
	"github.com/satori/go.uuid"
	"log"
	"mime/multipart"
	"time"
)

type VideoWithUser struct {
	repository.Video
	Author User `json:"author"`
}

func (videowithuser *VideoWithUser) GetVideosByLastTime(lastTime time.Time, userId int64) ([]VideoWithUser, time.Time, error) {
	videos := make([]VideoWithUser, 0, config.VideoCount)
	//get video by lastTime, videoCount is 20
	tableVideos, err := repository.QueryVideosByLastTime(lastTime)
	if err != nil {
		log.Printf("repository.QueryVideosByLastTime() error：%v", err)
		return nil, time.Time{}, err
	}
	//append usrId and videos to videoService struct
	err = videowithuser.copyVideos(&videos, &tableVideos, userId)
	if err != nil {
		log.Printf("repository.copyVideos() error：%v", err)
		return nil, time.Time{}, err
	}
	//return videos and first publishTime
	return videos, tableVideos[len(tableVideos)-1].PublishTime, nil
}

// Publish upload publish video to ftp
func (videowithuser *VideoWithUser) Publish(data *multipart.FileHeader, userId int64, title string) error {
	//upload publish video to sftp
	file, err := data.Open()
	if err != nil {
		log.Printf("videoService.data.Open() error: %v", err)
		return err
	}
	defer file.Close()

	// generate a uuid as videoname
	videoName := uuid.NewV4().String()

	err = repository.VideoSFTP(file, videoName)
	if err != nil {
		log.Printf("repository.VideoFTP() error%v", err)
		return err
	}
	//get screenshot
	imageName := uuid.NewV4().String()
	util.Ffchan <- util.Ffmsg{
		videoName,
		imageName,
	}
	//save data
	err = repository.SaveVideoToDB(videoName, imageName, userId, title)
	if err != nil {
		log.Printf("repository.Save(videoName, imageName, userId) error%v", err)
		return err
	}
	return nil
}

// List get videos by authorId
func (videowithuser *VideoWithUser) GetVideosByAuthorid(userId int64, curId int64) ([]VideoWithUser, error) {
	data, err := repository.QueryVideosByAuthorId(userId)
	if err != nil {
		log.Printf("repository.QueryVideosByAuthorId(%v) error:%v", userId, err)
		return nil, err
	}
	result := make([]VideoWithUser, 0, len(data))
	err = videowithuser.copyVideos(&result, &data, curId)
	if err != nil {
		log.Printf("videoService.copyVideos(&result, &data, %v) error:%v", userId, err)
		return nil, err
	}
	return result, nil
}

// copyVideos add curId and userid
func (videowithuser *VideoWithUser) copyVideos(result *[]VideoWithUser, data *[]repository.Video, userId int64) error {
	for _, temp := range *data {
		var video VideoWithUser
		//append userId and video
		video.Author, err = videowithuser.QueryUserByIdWithCurId(temp.AuthorId, userId)
		if err != nil {
			log.Printf("QueryUserByIdWithCurId(data.AuthorId, userId) error：%v", err)
		}
		*result = append(*result, video)
	}
	return nil
}
