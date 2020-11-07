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

// ユーザー登録関数
// request形式はbody json
func SignupPost(context *gin.Context) {
	var user models.UserModel
	if err := context.Bind(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err":err})
		return
	} else {
		err := CreateUser(user.Username, user.Password)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"err":err})
			return
		} else {
			context.JSON(http.StatusOK, gin.H{"message":"signup success"})
		}
	}
}
// パスワードをハッシュ化、ユーザー情報を保存する関数
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

// ログイン関数
// ログイン後にJWT Tokenを発行する
func LoginPost(context *gin.Context)  {
	var user models.UserModel
	if err := context.Bind(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err":err})
		return
	} else {
		db := config.DbConnect()
		loginUser := user.Username	// db query時に変数が入れ替わってしまうので、先に変数を定義しておきます。
		loginPassword := user.Password
		db.Find(&models.UserModel{}, "username =?", loginUser).Scan(&user)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPassword)); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"err":err, "login":loginUser})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"token": CreateJWTToken(user.Username),
		})
	}
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
