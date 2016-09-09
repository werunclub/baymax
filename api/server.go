package main

import (
	"baymax/api/handler/v1"
	mw "baymax/api/middleware"
	util "baymax/api/util"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func StartServer(addr string) {

	binding.Validator = &util.Validator{}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(mw.Recovery())
	//router.Use(gin.Recovery())
	//router.Use(mw.MaxAllowed(20))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "404", "message": "Endpoint not found"})
	})
	authMiddleware := mw.AuthMiddlewareInit("api auth", "secret key")

	api_v1 := router.Group("/v1")
	{
		userHandler := v1.NewUserHandler()
		api_v1.POST("/auth/login", authMiddleware.LoginHandler)

		api_v1.GET("/club/:clubId", v1.GetClub)
		api_v1.POST("/club/:clubId", v1.UpdateClub)
		api_v1.POST("/clubs", v1.CreateClub)

		{
			auth := api_v1.Group("/")
			auth.Use(authMiddleware.MiddlewareFunc())
			auth.GET("/auth/refresh_token", authMiddleware.RefreshHandler)
			auth.POST("/clubs", v1.CreateClub)
			auth.PATCH("/user/self", userHandler.PatchUser)
		}
	}

	endless.ListenAndServe(addr, router)
}
