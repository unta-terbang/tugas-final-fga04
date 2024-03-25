package controllers

import (
	"net/http"

	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"main.go/database"
	"main.go/models"
)



func CreateComment(ctx *gin.Context) {
    db := database.GetDB()

    userData := ctx.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))

    var comment models.Comment
    if err := ctx.ShouldBindJSON(&comment); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Gagal memvalidasi data",
            "error":   err.Error(),
        })
        return
    }

    var validationErrors []ValidationError

    if comment.Message == "" {
        validationErrors = append(validationErrors, ValidationError{Field: "message", Message: "Message harus diisi"})
    }

    photoIDStr := strconv.Itoa(int(comment.PhotoID))
    if comment.PhotoID == 0 {
        validationErrors = append(validationErrors, ValidationError{Field: "photo_id", Message: "Photo ID harus diisi"})
    } else if _, err := strconv.Atoi(photoIDStr); err != nil {
        validationErrors = append(validationErrors, ValidationError{Field: "photo_id", Message: "Photo ID harus berupa angka integer"})
    }


    if len(validationErrors) > 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Terjadi kesalahan dalam validasi",
            "errors":  validationErrors,
        })
        return
    }

    comment.UserID = userID

    if err := db.Create(&comment).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    response := gin.H{
        "id":       comment.Id,
        "message":  comment.Message,
        "photo_id": comment.PhotoID,
        "user_id": comment.UserID,
    }

    ctx.JSON(http.StatusCreated, response)
}


func GetAllUserComment(ctx *gin.Context) {
	db := database.GetDB()

	var comments []models.Comment
	if err := db.Find(&comments).Error; err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var commentsResponse []gin.H
	for _, comment := range comments {
		user := models.User{}
		db.First(&user, comment.UserID)

		var photo models.Photo
		if err := db.First(&photo, comment.PhotoID).Error; err != nil {
			// Jika tidak ada foto yang ditemukan, lanjutkan ke komentar berikutnya
			continue
		}

		commentResponse := gin.H{
			"id":         comment.Id,
			"message":    comment.Message,
			"user_id":    comment.UserID,
			"photo_id":   comment.PhotoID,
			"user": gin.H{
				"id":       user.Id,
				"email":    user.Email,
				"username": user.Username,
			},
			"photo": gin.H{
				"id":        photo.Id,
				"caption":   photo.Caption,
				"title":     photo.Title,
				"photo_url": photo.PhotoUrl,
			},
		}
		commentsResponse = append(commentsResponse, commentResponse)
	}

	ctx.JSON(http.StatusOK, commentsResponse)
}


func GetCommentByID(ctx *gin.Context) {
    db := database.GetDB()

    commentsID := ctx.Param("commentsId")

    id, err := strconv.Atoi(commentsID)
    if err != nil {
        ctx.AbortWithError(http.StatusBadRequest, errors.New("Invalid comments ID"))
        return
    }

    var comment models.Comment
    if err := db.Where("id = ?", id).First(&comment).Error; err != nil {
        ctx.AbortWithError(http.StatusNotFound, errors.New("Comment not found"))
        return
    }

    var user models.User
    if err := db.First(&user, comment.UserID).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    var photo models.Photo
    if err := db.First(&photo, comment.PhotoID).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    commentResponse := gin.H{
        "id":      comment.Id,
        "message": comment.Message,
        "user": gin.H{
            "id":       user.Id,
            "email":    user.Email,
            "username": user.Username,
        },
        "photo": gin.H{
            "id":        photo.Id,
            "caption":   photo.Caption,
            "title":     photo.Title,
            "photo_url": photo.PhotoUrl,
        },
        "user_id":  comment.UserID,
        "photo_id": comment.PhotoID,
    }

    ctx.JSON(http.StatusOK, commentResponse)
}

func UpdateCommentByID(ctx *gin.Context) {
    db := database.GetDB()

    commentID := ctx.Param("commentsId")

    var requestData map[string]interface{}
    if err := ctx.ShouldBindJSON(&requestData); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Gagal memvalidasi data",
            "error":   err.Error(),
        })
        return
    }

    message, ok := requestData["message"].(string)
    if !ok || message == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Message tidak boleh kosong",
        })
        return
    }

    var comment models.Comment
    if err := db.Where("id = ?", commentID).First(&comment).Error; err != nil {
        ctx.AbortWithError(http.StatusNotFound, errors.New("Comment tidak ditemukan"))
        return
    }

    updateData := make(map[string]interface{})
    updateData["message"] = message

    if err := db.Model(&comment).Updates(updateData).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    commentDataJSON := map[string]interface{}{
        "photo_id": comment.PhotoID,
        "user_id":  comment.UserID,
        "id":       comment.Id,
        "message":  updateData["message"],
    }

    ctx.JSON(http.StatusOK, commentDataJSON)
}


func DeleteCommentByID(ctx *gin.Context) {
    db := database.GetDB()

    commentID := ctx.Param("commentsId")

    var comment models.Comment
    if err := db.Where("id = ?", commentID).First(&comment).Error; err != nil {
        ctx.AbortWithError(http.StatusNotFound, errors.New("Foto tidak ditemukan"))
        return
    }

	if err := db.Delete(&comment).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Foto berhasil dihapus",
    })
}
