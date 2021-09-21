package server

import (
	"github.com/MISW/birdol-server/controller"
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
			user.DELETE("", controller.UnlinkAccount())
			account := user.Group("/account") //アカウント連携
			{
				account.POST("", controller.LinkAccount())
				account.PUT("/:userid", controller.SetDataLink())
			}
		}

		auth := v1.Group("/auth")
		{
			auth.POST("", controller.TokenAuthorize()) // Login using Token
		}
		progress := v1.Group("/progress/:userid")
		{
			progress.GET("/gallery", controller.GetGalleryInfo())
			progress.PUT("/complete", controller.FinishProgress())
			progress.GET("/complete", controller.GetCompletedCharacters())
			progress.GET("", controller.GetCurrentProgress())
			progress.PUT("", controller.CreateOrUpdateProgress())
		}
	}
	return router
}
