package main

import (
	"jwt-gin/api_v1/config"
	"jwt-gin/api_v1/models"
	"jwt-gin/api_v1/routes"

	"github.com/gin-gonic/gin"
)

func main()  {
	gin.SetMode(gin.DebugMode)
	db := config.DbConnect()
	defer db.Close()
	db.AutoMigrate(&models.UserModel{})
	router := routes.SetupRouter()
	router.Run(":8080")
}