package controllers

import (
	"jwt-gin/api_v1/config"
	"jwt-gin/api_v1/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignupGet(context *gin.Context)  {
	context.HTML(http.StatusOK, "signup.html", gin.H{
		"title":"Signup",
	})
}

// ユーザー登録
func SignupPost(context *gin.Context)  {
	var user models.UserModel
	if err := context.Bind(&user); err != nil {
		context.HTML(http.StatusBadRequest, "signup.html", gin.H{"err":err})
		context.Abort()
	} else {
		username := context.PostForm("username")
		password := context.PostForm("password")
		if err := creaetUser(username, password); err != nil {
			context.HTML(http.StatusBadRequest, "signup.html", gin.H{"err":err})
			return
		}
		context.Redirect(302, "/")
	}
}

// パスワードハッシュ関数で保存するまで
func creaetUser(username, password string) []error {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	db := config.DbConnect()
	defer db.Close()
	if err := db.Create(&models.UserModel{Username: username, Password: string(passwordHash)}).GetErrors(); err != nil {
		return err
	}
	return nil
}