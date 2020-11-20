package routes

import (
	"github.com/Watson-Sei/jwt-gin/api_v1/controllers"
	"github.com/Watson-Sei/jwt-gin/api_v1/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/api/signup", controllers.SignupPost)
		v1.POST("/api/login", controllers.LoginPost)
		v1.GET("/api/logout", middleware.JWTChecker(), controllers.LogoutGET)
		v1.GET("/api/refresh", middleware.RefreshChecker(), controllers.RefreshGET)
	}
	private := router.Group("/private")
	{
		private.GET("/book", middleware.JWTChecker(), controllers.BookGet)
	}
	return router
}