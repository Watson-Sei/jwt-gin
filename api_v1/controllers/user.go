package controllers

import (
	"jwt-gin/api_v1/config"
	"jwt-gin/api_v1/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignupGet(context *gin.Context)  {
	context.HTML(200, "signup.html", gin.H{})
}

func SignupPost(context *gin.Context) {
	var user models.UserModel
	if err := context.Bind(&user); err != nil {
		context.HTML(http.StatusBadRequest, "signup.html", gin.H{"err":err})
		context.Abort()
	} else {
		username := context.PostForm("username")
		password := context.PostForm("password")
		err := CreateUser(username, password)
		if err != nil {
			context.HTML(http.StatusBadRequest, "signup.html", gin.H{"err":err})
		}
		context.Redirect(302,"/v1/api/signup")
	}
}

func CreateUser(username, password string) (err error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	db := config.DbConnect()
	defer db.Close()
	if err := db.Create(&models.UserModel{Username: username, Password: string(hash)}).Error; err != nil {
		return err
	}
	return nil
}