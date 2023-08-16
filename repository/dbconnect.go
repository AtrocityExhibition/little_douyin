package repository

import (
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbConnection struct {
	db *gorm.DB
}

var (
	db_connection *dbConnection
	connectOnce   sync.Once
)

func GETInstanceDB() *dbConnection {
	connectOnce.Do(
		func() {
			db_connection = &dbConnection{}
			dsn := "root:root@tcp(127.0.0.1:3306)/VER_BASE?charset=utf8mb4&parseTime=True&loc=Local"
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
			if err != nil {
				panic("failed")
			}
			db_connection.db = db
		})
	return db_connection
}
