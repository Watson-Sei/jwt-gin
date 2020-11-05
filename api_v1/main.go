package main

import (
	. "jwt-gin/api_v1/config"
	. "jwt-gin/api_v1/models"
	. "jwt-gin/api_v1/routes"

	"github.com/gin-gonic/gin"
)

func main()  {
	gin.SetMode(gin.DebugMode)
	db := DbConnect()
	defer db.Close()
	db.AutoMigrate(&UserModel{})
	router := SetupRouter()
	router.Run(":8080")
}