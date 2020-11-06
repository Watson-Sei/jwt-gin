package middleware

import (
	"fmt"
	"jwt-gin/api_v1/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

// jwt検証
func JWTChecker() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := request.ParseFromRequest(context.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			b := []byte(config.SECRETKEY)
			return b, nil
		})

		if err == nil {
			claims := token.Claims.(jwt.MapClaims)
			msg := fmt.Sprintf("こんにちは、「%s」", claims["username"])
			context.JSON(200, gin.H{"message": msg})
		} else {
			context.JSON(401, gin.H{"err": fmt.Sprint(err)})
			context.Abort()
		}
	}
}