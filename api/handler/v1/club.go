package v1

import "github.com/gin-gonic/gin"

func GetClubDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": c.Param("clubId"),
	})
}

func CreateClub(c *gin.Context) {
	userId, ok := c.Get("userID")

	if ok == false {
		c.JSON(403, gin.H{"message": "bad access"})
		return
	}

	c.JSON(200, gin.H{
		"message": "create club",
		"userId":  userId,
	})
}
