package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"main.go/database"
	"main.go/helpers"
	"main.go/models"
)

func ReqAuth (ctx *gin.Context) {
     fmt.Println("TES MIDDLEWARE 1")
	 ctx.Next()
}


func Authentication() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		claim, err := helpers.VerifyToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("userData", claim)
		ctx.Next()
	}
}

func UserAuthorization() gin.HandlerFunc {
    return func(c *gin.Context) {
        db := database.GetDB()

        userData, exists := c.Get("userData")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "message": "Unauthorized",
                "error":   "User data not found in context",
            })
            return
        }
        userID := uint(userData.(jwt.MapClaims)["id"].(float64))

        var user models.User
        if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
                "message": "User not found",
                "error":   "User data not found in the database",
            })
            return
        }

        c.Set("user", user)

        c.Next()
    }
}



func PhotoAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		photoID, err := strconv.Atoi(c.Param("photoId"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Unauthorized",
				"error":   "Invalid photo ID data type",
			})
			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))
		photo := models.Photo{}

		err = db.Select("user_id").First(&photo, uint(photoID)).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unauthorized",
				"error":   "Failed to find photo",
			})
			return
		}

		if photo.UserID != userID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Forbidden",
				"error":   "You are not allowed to access this photo",
			})
			return
		}

		c.Next()
	}
}

func CommentAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		commentsID, err := strconv.Atoi(c.Param("commentsId"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Unauthorized",
				"error":   "Invalid Comment ID data type",
			})
			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))
		comment := models.Comment{}

		err = db.Select("user_id").First(&comment, uint(commentsID)).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unauthorized",
				"error":   "Failed to find Comment",
			})
			return
		}

		if comment.UserID != userID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Forbidden",
				"error":   "You are not allowed to access this Comment",
			})
			return
		}

		c.Next()
	}
}

func GetSocialMediaAuthorization() gin.HandlerFunc {
    return func(c *gin.Context) {
        db := database.GetDB()

        userData, exists := c.Get("userData")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "message": "Unauthorized",
                "error":   "User data not found in context",
            })
            return
        }
        userID := uint(userData.(jwt.MapClaims)["id"].(float64))

        var sosmed models.SocialMedia
        if err := db.Where("user_id = ?", userID).First(&sosmed).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
                "message": "Not Found",
                "error":   "Social media data not found for this user",
            })
            return
        }

        c.Next()
    }
}


func SocialMediaAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		sosmedID, err := strconv.Atoi(c.Param("socialmediaId"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Unauthorized",
				"error":   "Invalid social media ID data type",
			})
			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))
		sosmed := models.SocialMedia{}

		err = db.Select("user_id").First(&sosmed, uint(sosmedID)).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unauthorized",
				"error":   "Failed to find social media ID",
			})
			return
		}

		if sosmed.UserID != userID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Forbidden",
				"error":   "You are not allowed to access this social media ID",
			})
			return
		}

		c.Next()
	}
}
