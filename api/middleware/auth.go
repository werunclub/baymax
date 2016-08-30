package middleware

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"time"
)

func AuthMiddlewareInit(realm string, secretKey string) (authMiddleware *jwt.GinJWTMiddleware) {

	// the jwt middleware
	authMiddleware = &jwt.GinJWTMiddleware{
		Realm:      realm,
		Key:        []byte(secretKey),
		Timeout:    time.Hour * 24 * 365 * 2,
		MaxRefresh: time.Minute * 30,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if password == "test" {
				return userId, true
			}

			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	}

	return
}
