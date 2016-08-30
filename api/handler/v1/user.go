package v1

import "github.com/gin-gonic/gin"

type UserHandler struct {
}

func NewUserHandler() UserHandler {
	return UserHandler{}
}

func (handler *UserHandler) UserProfile(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
