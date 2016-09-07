package v1

import (
	"github.com/gin-gonic/gin"

	"baymax/errors"
)

type CreateClubRequest struct {
	// 城市ID
	CityCode string `form:"city_code" json:"city_code" xml:"city_code"`
	// 俱乐部描述
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	// 行业ID
	IndustryID int `json:"industry_id" valid:"min=21"`
	// 俱乐部名称
	Name string `form:"name" json:"name" xml:"name"`
}

func GetClubDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": c.Param("clubId"),
	})
}

func CreateClub(c *gin.Context) {
	//userId, ok := c.Get("userID")

	var req CreateClubRequest

	err := c.Bind(&req)

	if err != nil {
		c.AbortWithError(400, err)
		c.JSON(400, errors.BadRequest(err.Error()))
		return
	}

	c.JSON(200, req)
}
