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
v1: For development	*/

func GetRouterV1(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	userController := controller.UserController{DB: db, TokenManager: &auth.TokenManager{DB: db}}
	authController := controller.AuthController{DB: db, TokenManager: &auth.TokenManager{DB: db}, SessionManager: &auth.SessionManager{DB: db}}
	securityMiddleware := security.SecurityMiddleware{DB: db}
	sessionMiddleware := session.SessionMiddleware{DB: db}

	v1 := router.Group("api/v1")
	{
		user := v1.Group("/user")
		{
			user.PUT("", userController.HandleSignUp())
			user.POST("", userController.LinkAccount())
		}

		auth := v1.Group("/session")
		auth.Use(securityMiddleware.RequestValidation())
		{
			auth_root := auth.Group("")
			auth_root.Use(sessionMiddleware.CheckToken())
			{
				auth_root.GET("", authController.TokenAuthorize())   // Login using Token
				auth_root.DELETE("", authController.UnlinkAccount()) // Unlink Account
				auth_root.PUT("", authController.SetDataLink())      // Link Account
			}
			auth.GET("/refresh", authController.RefreshToken()) // Token Refresh
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
