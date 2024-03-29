package server

import (
	"github.com/MISW/birdol-server/controller"
	"github.com/MISW/birdol-server/middlewares"
	"github.com/gin-gonic/gin"
)

/*	ルーティングはここで設定する
v2: For production	*/

func GetRouterV2() *gin.Engine {
	router := gin.Default()
	v2 := router.Group("api/v2")
	{
		cli := v2.Group("/cli")
		{
			cli.POST("/version", controller.ClientVerCheck())
		}
		
		user := v2.Group("/user")
		{
			user.PUT("", controller.HandleSignUp())
			user.POST("", controller.LinkAccount())
		}

		auth := v2.Group("/auth")
		auth.Use(middlewares.RequestValidation())
		{
			auth_root := auth.Group("")
			auth_root.Use(middlewares.CheckToken())
			{
				auth_root.GET("", controller.TokenAuthorize()) // Login using Token
				auth_root.DELETE("", controller.UnlinkAccount()) // Unlink Account
				auth_root.PUT("", controller.SetDataLink()) // Link Account
			}
		}
		
		refresh := v2.Group("/refresh")
		refresh.Use(middlewares.RequestValidation())
		{
			refresh.GET("", controller.RefreshToken()) // Token Refresh
		}

		gamedata := v2.Group("/gamedata")
		gamedata.Use(middlewares.RequestValidation())
		gamedata.Use(middlewares.CheckToken())
		{
			// UNIMPLEMENTED
			gamedata_nobody := gamedata.Group("")
			gamedata_nobody.Use(middlewares.ReadSessionIDfromQuery())
			{
				// Requests that have no body
				gamedata_nobody.GET("/gallery", controller.GetGalleryInfo())
				gamedata_nobody.GET("/complete", controller.GetCompletedCharacters())
				gamedata_nobody.GET("/character", controller.GetCurrentCharacters())
				gamedata_nobody.GET("/story", controller.GetCurrentStory())
				
			}
			gamedata_body := gamedata.Group("")
			gamedata_body.Use(middlewares.ReadSessionIDfromBody())
			{
				// Requests that have body
				gamedata_body.PUT("/complete", controller.FinishProgress())
				gamedata_body.PUT("/character", controller.UpdateCharacters())
				gamedata_body.PUT("/story", controller.UpdateMainStory())
				gamedata_body.PUT("/new", controller.CreateProgress())
			}
		}
	}
	return router
}
