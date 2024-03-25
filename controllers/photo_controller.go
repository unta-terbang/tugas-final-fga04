package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"main.go/database"
	"main.go/models"
)


func CreatePhoto(ctx *gin.Context) {
    db := database.GetDB()

    userData := ctx.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))

    var photo models.Photo
    if err := ctx.ShouldBindJSON(&photo); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Gagal memvalidasi data",
            "error":   err.Error(),
        })
        return
    }

    var validationErrors []ValidationError

    if photo.Title == "" {
        validationErrors = append(validationErrors, ValidationError{Field: "title", Message: "Title harus diisi"})
    }

    if photo.PhotoUrl == "" {
        validationErrors = append(validationErrors, ValidationError{Field: "photo_url", Message: "Photo URL harus diisi"})
    } else if !strings.HasPrefix(photo.PhotoUrl, "http://") && !strings.HasPrefix(photo.PhotoUrl, "https://") {
        validationErrors = append(validationErrors, ValidationError{Field: "photo_url", Message: "Photo URL harus diawali dengan http:// atau https://"})
    }

    if len(validationErrors) > 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Terjadi kesalahan dalam validasi",
            "errors":  validationErrors,
        })
        return
    }

    photo.UserID = userID

    if err := db.Create(&photo).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    response := gin.H{
        "id":        photo.Id,
        "caption":   photo.Caption,
        "title":     photo.Title,
        "user_id":   photo.UserID,
        "photo_url": photo.PhotoUrl,
    }

    ctx.JSON(http.StatusCreated, response)
}


func GetAllUserPhotos(ctx *gin.Context) {
    db := database.GetDB()

    var photos []models.Photo
    if err := db.Find(&photos).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    var photosResponse []gin.H
    for _, photo := range photos {
        user := models.User{}
        db.First(&user, photo.UserID)

        photoResponse := gin.H{
            "id":        photo.Id,
            "caption":   photo.Caption,
            "title":     photo.Title,
            "photo_url": photo.PhotoUrl,
            "user_id":   photo.UserID,
            "user": gin.H{
                "id":       user.Id,
                "email":    user.Email,
                "username": user.Username,
            },
        }
        photosResponse = append(photosResponse, photoResponse)
    }

    ctx.JSON(http.StatusOK, photosResponse)
}


func GetPhotoByID(ctx *gin.Context) {
    db := database.GetDB()

    photoIDStr := ctx.Param("photoId")
    if photoIDStr == "" {
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "message": "Photo ID is required",
        })
        return
    }

    photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
    if err != nil {
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "message": "Invalid Photo ID",
        })
        return
    }

    var photo models.Photo
    if err := db.Where("id = ?", photoID).First(&photo).Error; err != nil {
        ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
            "message": "Photo not found",
        })
        return
    }

    var user models.User
    if err := db.Where("id = ?", photo.UserID).First(&user).Error; err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to find user",
        })
        return
    }

    response := gin.H{
        "id":        photo.Id,
        "caption":   photo.Caption,
        "title":     photo.Title,
        "photo_url": photo.PhotoUrl,
        "user_id":   photo.UserID,
        "user": gin.H{
            "id":       user.Id,
            "email":    user.Email,
            "username": user.Username,
        },
    }

    ctx.JSON(http.StatusOK, response)
}

func UpdatePhotoByID(ctx *gin.Context) {
    db := database.GetDB()

    photoID := ctx.Param("photoId")

    var requestData map[string]interface{}
    if err := ctx.ShouldBindJSON(&requestData); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Gagal memvalidasi data",
            "error":   err.Error(),
        })
        return
    }

    var photo models.Photo
    if err := db.Where("id = ?", photoID).First(&photo).Error; err != nil {
        ctx.AbortWithError(http.StatusNotFound, errors.New("Photo tidak ditemukan"))
        return
    }

    updateData := make(map[string]interface{})

    if title, ok := requestData["title"].(string); ok {
        updateData["title"] = title
    }
    if caption, ok := requestData["caption"].(string); ok {
        updateData["caption"] = caption
    }
    if photoURL, ok := requestData["photo_url"].(string); ok {
        updateData["photo_url"] = photoURL
    }

    var validationErrors []ValidationError

    if title, exists := updateData["title"].(string); !exists || title == "" {
        validationErrors = append(validationErrors, ValidationError{Field: "title", Message: "Title harus diisi"})
    }

    if photoURL, exists := updateData["photo_url"].(string); !exists || photoURL == "" {
        validationErrors = append(validationErrors, ValidationError{Field: "photo_url", Message: "Photo URL harus diisi"})
    } else if !strings.HasPrefix(photoURL, "http://") && !strings.HasPrefix(photoURL, "https://") && !strings.HasPrefix(photoURL, "https://") {
        validationErrors = append(validationErrors, ValidationError{Field: "photo_url", Message: "Photo URL harus diawali dengan http:// atau https://"})
    }

    if len(validationErrors) > 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Terjadi kesalahan dalam validasi",
            "errors":  validationErrors,
        })
        return
    }

    if err := db.Model(&photo).Updates(updateData).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    photoDataJSON := map[string]interface{}{
        "id":         photo.Id,
        "title":      updateData["title"],
        "caption":    updateData["caption"],
        "photo_url":  updateData["photo_url"],
        "user_id":    photo.UserID,
    }

    ctx.JSON(http.StatusOK, photoDataJSON)
}



func DeletePhotoByID(ctx *gin.Context) {
    db := database.GetDB()

    photoID := ctx.Param("photoId")

    var photo models.Photo
    if err := db.Where("id = ?", photoID).First(&photo).Error; err != nil {
        ctx.AbortWithError(http.StatusNotFound, errors.New("Foto tidak ditemukan"))
        return
    }

    if err := db.Delete(&photo).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Foto berhasil dihapus",
    })
}
