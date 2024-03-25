package router

import (
	"github.com/gin-gonic/gin"
	"main.go/controllers"
	"main.go/middleware"
)

func Routers() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("users")
	{
		userRouter.POST("/register", controllers.CreateUser)
		userRouter.POST("/login", controllers.LoginUser)
		userRouter.GET("/validates", middleware.Authentication(), controllers.ValidateTes)
		userRouter.PUT("", middleware.Authentication(),middleware.UserAuthorization(), controllers.UpdateUser)
		userRouter.DELETE("", middleware.Authentication(),middleware.UserAuthorization(), controllers.DeleteUser)
	}

	photoRouter := r.Group("photos")
	{
		photoRouter.Use(middleware.Authentication())
		photoRouter.POST("", controllers.CreatePhoto)
		photoRouter.GET("", controllers.GetAllUserPhotos)	
		photoRouter.GET("/:photoId", controllers.GetPhotoByID)
		photoRouter.PUT("/:photoId", middleware.PhotoAuthorization() ,controllers.UpdatePhotoByID)
		photoRouter.DELETE("/:photoId", middleware.PhotoAuthorization() ,controllers.DeletePhotoByID)
	}

	commentRouter := r.Group("comments")
	{
		commentRouter.Use(middleware.Authentication())
		commentRouter.POST("", controllers.CreateComment)
		commentRouter.GET("", controllers.GetAllUserComment)
		commentRouter.GET("/:commentsId", controllers.GetCommentByID)
		commentRouter.PUT("/:commentsId",middleware.CommentAuthorization(), controllers.UpdateCommentByID)
		commentRouter.DELETE("/:commentsId",middleware.CommentAuthorization(), controllers.DeleteCommentByID)
	}
	
	socialMediaRouter := r.Group("socialmedias")
	{
		socialMediaRouter.Use(middleware.Authentication())
		socialMediaRouter.POST("", controllers.CreateSocialMedia)
		socialMediaRouter.GET("",middleware.GetSocialMediaAuthorization(), controllers.GetUserSocialMedia)
		socialMediaRouter.GET("/:socialmediaId",middleware.SocialMediaAuthorization(), controllers.GetSocialMediaByID)
		socialMediaRouter.PUT("/:socialmediaId",middleware.SocialMediaAuthorization(), controllers.UpdateSocialMediaByID)
		socialMediaRouter.DELETE("/:socialmediaId",middleware.SocialMediaAuthorization(), controllers.DeleteSocialMediaByID)
	}

	return r
}
