package repository

import (
	"time"
)

type Video struct {
	Id         int       `gorm:"primary_key"`
	Author_id  int       `gorm:"not null"`
	Title      string    `gorm:"not null"`
	Createtime time.Time `gorm:"autoCreateTime;type:datetime"`
}

type User struct {
	Id       int    `gorm:"primary_key"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// 查询video /douyin/feed 顶多显示30个结果
func (m *dbConnection) SearchVideosbyTime(t time.Time) ([]Video, error) {
	videos := make([]Video, 0)
	res := m.db.Where("createtime <= ?", t).Order("createtime desc").Limit(30).Find(&videos)
	return videos, res.Error
}

// 插入User /douyin/user/register
func (m *dbConnection) InsertUser(username string, password string) (User, error) {
	u := User{}
	u.Username = username
	u.Password = password
	res := m.db.Create(&u)
	return u, res.Error
}

// 登录User /douyin/user/login
func (m *dbConnection) LoginUser(username string, password string) (User, error) {
	u := User{}
	res := m.db.Where("username = ?", username).First(&u)
	if res.Error != nil {
		return u, res.Error
	}
	if u.Password != password {
		u.Id = -1
	}
	return u, nil
}

// 查询User, 传出参数 user, /douyin/user
func (m *dbConnection) ShowUser(uid int) (User, error) {
	u := User{}
	res := m.db.Where("id = ?", uid).First(&u)
	return u, res.Error
}

// 插入video /douyin/publish/action/
func (m *dbConnection) InsertVideo(v Video) error {
	res := m.db.Create(&v)
	return res.Error
}

// 查询用户pub videos
func (m *dbConnection) SearchPublished(aid int) ([]Video, error) {
	pubvideos := make([]Video, 0)
	res := m.db.Where("author_id = ?", aid).Find(&pubvideos)
	return pubvideos, res.Error
}
