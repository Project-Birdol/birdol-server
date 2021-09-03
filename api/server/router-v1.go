package server
import (
	"github.com/gin-gonic/gin"
	"github.com/MISW/birdol-server/controller"
)

/*	ルーティングはここで設定する
		  v1: For development	*/

func GetRouterV1() *gin.Engine{
	router := gin.Default()
	v1 := router.Group("api/v1")
	{
		user := v1.Group("/user")
		{
			user.PUT("", controller.HandleSignUp())
			user.POST("", controller.HandleLogin())
			user.DELETE("", controller.HandleLogout())
		}

		auth := v1.Group("/auth")
		{
			auth.POST("", controller.TokenAuthorize()) // Login using Token
		}
	}
	return router
}