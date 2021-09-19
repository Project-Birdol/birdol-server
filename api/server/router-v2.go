package server

import (
	"github.com/MISW/birdol-server/controller"
	"github.com/gin-gonic/gin"
)

/*	ルーティングはここで設定する
v2: For production	*/

func GetRouterV2() *gin.Engine {
	router := gin.Default()
	v2 := router.Group("api/v2")
	{
		user := v2.Group("/user")
		{
			user.PUT("", controller.HandleSignUp())
			user.DELETE("", controller.HandleLogout())
			account := user.Group("/account") //アカウント連携
			{
				account.POST("", controller.LoginAccount())
				account.PUT("/:userid", controller.EditAccount())
			}
		}

		auth := v2.Group("/auth")
		{
			auth.POST("", controller.TokenAuthorize()) // Login using Token
		}
	}
	return router
}
