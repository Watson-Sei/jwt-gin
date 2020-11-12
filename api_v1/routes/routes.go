package routes

import (
	"jwt-gin/api_v1/controllers"
	"jwt-gin/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/api/signup", controllers.SignupPost)
		v1.POST("/api/login", controllers.LoginPost)
		v1.POST("/api/logout", middleware.JWTChecker(), controllers.LogoutPost)
	}
	private := router.Group("/private")
	{
		private.GET("/book", middleware.JWTChecker(), controllers.BookGet)
	}
	return router
}