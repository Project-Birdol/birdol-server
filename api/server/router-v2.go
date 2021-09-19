package server
import (
	"github.com/gin-gonic/gin"
	"github.com/MISW/birdol-server/controller"
	"github.com/MISW/birdol-server/middlewares"
)

/*	ルーティングはここで設定する
		  v2: For production	*/

func GetRouterV2() *gin.Engine{
	router := gin.Default()
	v2 := router.Group("api/v2")
	{
		user := v2.Group("/user")
		{
			user.PUT("", controller.HandleSignUp())
			user.POST("", controller.HandleLogin())
			user.DELETE("", controller.HandleLogout())
		}

		auth := v2.Group("/auth")
		auth.Use(middlewares.RequestValidation())
		{
			auth.POST("", controller.TokenAuthorize()) // Login using Token
		}
	}
	return router
}