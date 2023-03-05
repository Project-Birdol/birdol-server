package server

import (
	"github.com/Project-Birdol/birdol-server/auth"
	"github.com/Project-Birdol/birdol-server/controller"
	security "github.com/Project-Birdol/birdol-server/middleware/security"
	session "github.com/Project-Birdol/birdol-server/middleware/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*	ルーティングはここで設定する
v2: For production	*/

func GetRouterV2(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	userController := controller.UserController{DB: db, TokenManager: &auth.TokenManager{DB: db}}
	progressController := controller.ProgressController{DB: db}
	versionController := controller.VersionController{DB: db}
	authController := controller.AuthController{DB: db, TokenManager: &auth.TokenManager{DB: db}, SessionManager: &auth.SessionManager{DB: db}}
	securityMiddleware := security.SecurityMiddleware{DB: db}
	sessionMiddleware := session.SessionMiddleware{DB: db}
	v2 := router.Group("api/v2")
	{
		cli := v2.Group("/cli")
		{
			cli.POST("/version", versionController.ClientVerCheck())
		}

		user := v2.Group("/user")
		user.Use(securityMiddleware.InspectPublicKey())
		{
			user.PUT("", userController.HandleSignUp())
			user.POST("", userController.LinkAccount())
		}

		auth := v2.Group("/auth")
		auth.Use(securityMiddleware.RequestValidation())
		{
			auth_root := auth.Group("")
			auth_root.Use(sessionMiddleware.CheckToken())
			{
				auth_root.GET("", authController.TokenAuthorize())   // Login using Token
				auth_root.DELETE("", authController.UnlinkAccount()) // Unlink Account
				auth_root.PUT("", authController.SetDataLink())      // Link Account
			}
		}

		refresh := v2.Group("/refresh")
		refresh.Use(securityMiddleware.RequestValidation())
		{
			refresh.GET("", authController.RefreshToken()) // Token Refresh
		}

		gamedata := v2.Group("/gamedata")
		gamedata.Use(securityMiddleware.RequestValidation())
		gamedata.Use(sessionMiddleware.CheckToken())
		{
			// UNIMPLEMENTED
			gamedata_nobody := gamedata.Group("")
			gamedata_nobody.Use(sessionMiddleware.ReadSessionIDfromQuery())
			{
				// Requests that have no body
				gamedata_nobody.GET("/gallery", progressController.GetGalleryInfo())
				gamedata_nobody.GET("/complete", progressController.GetCompletedCharacters())
				gamedata_nobody.GET("/character", progressController.GetCurrentCharacters())
				gamedata_nobody.GET("/story", progressController.GetCurrentStory())

			}
			gamedata_body := gamedata.Group("")
			gamedata_body.Use(sessionMiddleware.ReadSessionIDfromBody())
			{
				// Requests that have body
				gamedata_body.PUT("/complete", progressController.FinishProgress())
				gamedata_body.PUT("/character", progressController.UpdateCharacters())
				gamedata_body.PUT("/story", progressController.UpdateMainStory())
				gamedata_body.PUT("/new", progressController.CreateProgress())
			}
		}
	}
	return router
}
