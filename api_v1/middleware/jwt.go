package middleware

import (
	"github.com/Watson-Sei/jwt-gin/api_v1/config"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

// jwt検証
func JWTChecker() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := request.ParseFromRequest(context.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			b := []byte(config.ACCESS_TOKEN_SECRETKEY)
			return b, nil
		})

		if err == nil {
			err := BlackListChecker(token.Raw)
			if err == nil {
				context.JSON(http.StatusBadRequest, gin.H{"err":"無効トークン"})
				context.Abort()
			} else {
				claims := token.Claims.(jwt.MapClaims)
				context.Set("userId", claims["userId"])
				context.Set("username", claims["username"])
				context.Set("exp", claims["exp"])
				context.Set("token", token.Raw)
				context.Next()
			}
		} else {
			context.JSON(http.StatusUnauthorized, gin.H{"err":err})
			context.Abort()
		}
	}
}

func RefreshChecker() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := request.ParseFromRequest(context.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			b := []byte(config.REFRESH_TOKEN_SECRET_KEY)
			return b, nil
		})

		if err == nil {
			err := BlackListChecker(token.Raw)
			if err == nil {
				context.JSON(http.StatusBadRequest, gin.H{"err":"無効リフレッシュトークンです"})
				context.Abort()
			} else {
				claims := token.Claims.(jwt.MapClaims)
				context.Set("userId", claims["userId"])
				context.Set("exp", claims["exp"])
				context.Set("token", token.Raw)
				context.Next()
			}
		} else {
			context.JSON(http.StatusUnauthorized, gin.H{"err":err})
			context.Abort()
		}
	}
}

// BlackListChecker
func BlackListChecker(token string) error {
	conn := config.RedisConnection()
	defer conn.Close()
	_, err := redis.String(conn.Do("GET", token))
	if err != nil {
		return err
	}
	return nil
}