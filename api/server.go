package main

import (
	"baymax/api/handler/v1"
	mw "baymax/api/middleware"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func StartServer(addr string) {

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	//router.Use(mw.MaxAllowed(20))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "404", "message": "Endpoint not found"})
	})
	authMiddleware := mw.AuthMiddlewareInit("api auth", "secret key")

	api_v1 := router.Group("/v1")
	{
		api_v1.POST("/auth/login", authMiddleware.LoginHandler)

		api_v1.GET("/club/:clubId", v1.GetClubDetail)

		userHandler := v1.NewUserHandler()

		auth := api_v1.Group("/")
		auth.Use(authMiddleware.MiddlewareFunc())
		{
			auth.GET("/auth/refresh_token", authMiddleware.RefreshHandler)

			auth.GET("/user", userHandler.UserProfile)

			auth.POST("/clubs", v1.CreateClub)
		}
	}

	endless.ListenAndServe(addr, router)
}
