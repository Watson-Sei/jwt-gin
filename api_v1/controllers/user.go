package controllers

import (
	"jwt-gin/api_v1/config"
	"jwt-gin/api_v1/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ユーザー登録関係
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
			return
		}
		context.Redirect(302,"/v1/api/signup")
	}
}

func CreateUser(username, password string) (err error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return err
	}
	db := config.DbConnect()
	defer db.Close()
	if err := db.Create(&models.UserModel{Username: username, Password: string(hash)}).Error; err != nil {
		return err
	}
	return nil
}

// ログイン
func LoginGet(context *gin.Context)  {
	context.HTML(http.StatusOK, "login.html", gin.H{"message":"まだログインしていません"})
}
// ログイン後にJWT Tokenを発行する
func LoginPost(context *gin.Context)  {
	db := config.DbConnect()
	var user models.UserModel
	db.Find(&models.UserModel{}, "username =?", context.PostForm("username")).Scan(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(context.PostForm("password"))); err != nil {
		context.HTML(http.StatusBadRequest, "signup.html", gin.H{"err":err})
		context.Abort()
	}
	context.JSON(http.StatusOK, gin.H{
		"token": CreateJWTToken(user.Username),
	})
}

// jwt token 生成
func CreateJWTToken(username string) string {
	/*
		アルゴリズムの指定
	*/
	// headerのセット
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 4).Unix()

	// 電子署名
	tokenString, err := token.SignedString([]byte(config.SECRETKEY))
	if err == nil {
		return tokenString
	} else {
		return "token生成に失敗しました。"
	}
}


// プライベートデータ

type Book struct {
	Title	string	`json:"title"`
	Tag 	string	`json:"tag"`
	URL 	string	`json:"url"`
}

func BookGet(context *gin.Context)  {
	context.JSON(http.StatusOK, Book{Title: "JWT認証を学ぼう！", Tag: "Golang", URL: "http://google.com"})
}
