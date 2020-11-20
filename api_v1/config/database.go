package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB gorm.DB
var count = 0
var err error
var ACCESS_TOKEN_SECRETKEY = os.Getenv("ACCESS_TOKEN_SECRETKEY")
var REFRESH_TOKEN_SECRET_KEY = os.Getenv("REFRESH_TOKEN_SECRETKEY")

type DBConfig struct {
	Host 		string
	User 		string
	DBName		string
	Password	string
}

func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
		Host:		"tcp(db)",
		User:		os.Getenv("MYSQL_USER"),
		Password: 	os.Getenv("MYSQL_PASSWORD"),
		DBName: 	os.Getenv("MYSQL_DATABASE"),
	}
	return &dbConfig
}

func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.DBName,
	)
}

func DbConnect() *gorm.DB {
	db, err := gorm.Open("mysql", DbURL(BuildDBConfig()))
	if err != nil {
		log.Println("Not ready, Retry connecting...")
		time.Sleep(time.Second)
		count++
		log.Println(count)
		if count > 30 {
			panic(err)
		}
		return DbConnect()
	}
	return db
}