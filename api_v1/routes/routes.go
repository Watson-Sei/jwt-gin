package routes

import (
	"jwt-gin/api_v1/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	v1 := router.Group("/v1")
	{
		v1.GET("/api/signup", controllers.SignupGet)
		v1.POST("/api/signup", controllers.SignupPost)
	}
	return router
}