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
		progress := v1.Group("/progress")
		{
			story := progress.Group("/story")
			{
				story.PUT("/update/:userid",controller.UpdateStory())
			}
			characters := progress.Group("/characters")
			{
				characters.GET("/gallary/:userid",controller.GetGallaryInfo())
				characters.GET("/completed/:userid",controller.GetCompletedCharacters())
				characters.PUT("/update/:userid",controller.UpdateCharacters())
			}
			progress.GET("/current/:userid",controller.GetCurrentProgress())
			progress.PUT("/new/:userid",controller.NewProgress())
		}
	}
	return router
}