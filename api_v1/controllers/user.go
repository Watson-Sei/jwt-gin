package controllers

import (
	"github.com/Watson-Sei/jwt-gin/api_v1/config"
	"github.com/Watson-Sei/jwt-gin/api_v1/models"
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
		tokens, err := CreateJWTToken(user.Username, user.ID)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"err": err})
			context.Abort()
		}
		context.JSON(http.StatusOK, tokens)
	}
}

// jwt token 生成
func CreateJWTToken(username string, userId uint) (map[string]string, error) {
	/*
		アルゴリズムの指定
	*/
	// headerのセット
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "AccessToken"
	claims["userId"] = userId
	claims["username"] = username
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 4).Unix()

	t, err := token.SignedString([]byte(config.ACCESS_TOKEN_SECRETKEY))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = "RefreshToken"
	rtClaims["userId"] = userId
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rt, err := refreshToken.SignedString([]byte(config.REFRESH_TOKEN_SECRET_KEY))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token": t,
		"refresh_token": rt,
	}, nil
}

// Logout
func LogoutGET(context *gin.Context)  {
	exp := context.MustGet("exp").(float64)
	token := context.MustGet("token").(string)
	err := BlackListSet(int64(exp), token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err":err})
		context.Abort()
	} else {
		context.JSON(http.StatusOK, gin.H{"message":"Logout成功しました"})
	}
}

// Redis JWT Token Black List Register
func BlackListSet(exp int64, accessToken string) error {
	conn := config.RedisConnection()
	defer conn.Close()

	// 残り時間
	nowTime := time.Now()
	expTime := time.Unix(exp, 0)

	// 残り時間秒数
	timeLeft := expTime.Sub(nowTime).Seconds()

	// Redis DBに追加
	_, err := conn.Do("SET", accessToken, string(exp))
	_, err = conn.Do("EXPIRE", accessToken, int64(timeLeft))
	if err != nil {
		return err
	}
	return nil
}

// Refresh
func RefreshGET(context *gin.Context)  {
	var user models.UserModel
	userId := context.MustGet("userId").(float64)
	db := config.DbConnect()
	defer db.Close()
	db.First(&user, uint(userId)).Scan(&user)
	// 新しいトークンを生成します
	tokens, err := CreateJWTToken(user.Username,uint(userId))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err":err})
		return
	} else {
		// リフレッシュトークンをブラックリストに保存します。
		exp := context.MustGet("exp").(float64)
		token := context.MustGet("token").(string)
		err := BlackListSet(int64(exp), token)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"err":err})
			context.Abort()
		} else {
			context.JSON(http.StatusOK, tokens)
		}
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
