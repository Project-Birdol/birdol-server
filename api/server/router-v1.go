package server

import (
	"github.com/MISW/birdol-server/controller"
	"github.com/MISW/birdol-server/middlewares"
	"github.com/gin-gonic/gin"
)

/*	ルーティングはここで設定する
v1: For development	*/

func GetRouterV1() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("api/v1")
	{
		user := v1.Group("/user")
		{
			user.PUT("", controller.HandleSignUp())
			user.POST("", controller.LinkAccount())
		}

		auth := v1.Group("/auth")
		auth.Use(middlewares.RequestValidation())
		{
			auth_root := auth.Group("")
			auth_root.Use(middlewares.CheckToken())
			{
				auth_root.GET("", controller.TokenAuthorize()) // Login using Token
				auth_root.DELETE("", controller.UnlinkAccount()) // Unlink Account
				auth_root.PUT("", controller.SetDataLink()) // Link Account
			}
			auth.GET("/refresh", controller.RefreshToken()) // Token Refresh
		}
		/*
		progress := v1.Group("/gamedata/:userid")
		{
			progress.GET("/gallery", controller.GetGalleryInfo())
			progress.PUT("/complete", controller.FinishProgress())
			progress.GET("/complete", controller.GetCompletedCharacters())
			progress.GET("", controller.GetCurrentProgress())
			progress.PUT("", controller.CreateOrUpdateProgress())
		}
		*/
	}
	return router
}
