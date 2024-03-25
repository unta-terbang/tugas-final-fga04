package controllers

import (
	"net/http"

	"errors"
	"strconv"
    "strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"main.go/database"
	"main.go/models"
)

func CreateSocialMedia(ctx *gin.Context) {
    db := database.GetDB()

    userData := ctx.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))

    var socialMedia models.SocialMedia
    if err := ctx.ShouldBindJSON(&socialMedia); err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    var validationErrors []ValidationError

    if socialMedia.Name == "" {
        validationErrors = append(validationErrors, ValidationError{Field: "name", Message: "Name harus diisi"})
    }

if socialMedia.SocialMediaUrl == "" {
    validationErrors = append(validationErrors, ValidationError{Field: "social_media_url", Message: "Social Media URL harus diisi"})
} else if !strings.HasPrefix(socialMedia.SocialMediaUrl, "http://") && !strings.HasPrefix(socialMedia.SocialMediaUrl, "https://") {
    validationErrors = append(validationErrors, ValidationError{Field: "social_media_url", Message: "Social Media URL harus diawali dengan http:// atau https://"})
}

    if len(validationErrors) > 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Terjadi kesalahan dalam validasi",
            "errors":  validationErrors,
        })
        return
    }

    socialMedia.UserID = userID

    if err := db.Create(&socialMedia).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    response := gin.H{
        "id":              socialMedia.Id,
        "name":            socialMedia.Name,
        "social_media_url": socialMedia.SocialMediaUrl,
        "user_id":         socialMedia.UserID,
    }

    ctx.JSON(http.StatusCreated, response)
}



func GetUserSocialMedia(ctx *gin.Context) {
    db := database.GetDB()
    userData := ctx.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))

    var socialMedia []models.SocialMedia
    if err := db.Where("user_id = ?", userID).Find(&socialMedia).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    var sosmedResponse []gin.H
    for _, sosmed := range socialMedia {
        user := models.User{}
        db.First(&user, sosmed.UserID)

        sosmedsResponse := gin.H{
            "id":        sosmed.Id,
            "name":   sosmed.Name,
            "social_media_url":     sosmed.SocialMediaUrl,
            "user": gin.H{
                "id":       user.Id,
                "email":    user.Email,
                "username": user.Username,
            },
        }
        sosmedResponse = append(sosmedResponse, sosmedsResponse)
    }

    ctx.JSON(http.StatusOK, sosmedResponse)
}


func GetSocialMediaByID(ctx *gin.Context) {
    db := database.GetDB()

    sosmedID := ctx.Param("socialmediaId")

    id, err := strconv.Atoi(sosmedID)
    if err != nil {
        ctx.AbortWithError(http.StatusBadRequest, errors.New("Invalid social media ID"))
        return
    }

    var sosmed models.SocialMedia
    if err := db.Where("id = ?", id).First(&sosmed).Error; err != nil {
        ctx.AbortWithError(http.StatusNotFound, errors.New("Comment not found"))
        return
    }

    var user models.User
    if err := db.First(&user, sosmed.UserID).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    sosmedResponse := gin.H{
        "id":      sosmed.Id,
        "name": sosmed.Name,
		"social_media_url":  sosmed.SocialMediaUrl,
        "user": gin.H{
            "id":       user.Id,
            "email":    user.Email,
            "username": user.Username,
        },
    }

    ctx.JSON(http.StatusOK, sosmedResponse)
}

func UpdateSocialMediaByID(ctx *gin.Context) {
    db := database.GetDB()

    sosmedID := ctx.Param("socialmediaId")

    var requestData map[string]interface{}
    if err := ctx.ShouldBindJSON(&requestData); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Gagal memvalidasi data",
            "error":   err.Error(),
        })
        return
    }

    var validationErrors []ValidationError

    if name, ok := requestData["name"].(string); ok {
        if name == "" {
            validationErrors = append(validationErrors, ValidationError{Field: "name", Message: "Name harus diisi"})
        }
    } else {
        validationErrors = append(validationErrors, ValidationError{Field: "name", Message: "Name harus berupa string"})
    }

    if SocialMediaUrl, ok := requestData["social_media_url"].(string); ok {
        if SocialMediaUrl == "" {
            validationErrors = append(validationErrors, ValidationError{Field: "social_media_url", Message: "Social Media URL harus diisi"})
        } else if !strings.HasPrefix(SocialMediaUrl, "http://") && !strings.HasPrefix(SocialMediaUrl, "https://") {
            validationErrors = append(validationErrors, ValidationError{Field: "social_media_url", Message: "Social Media URL harus diawali dengan http:// atau https://"})
        }
    } else {
        validationErrors = append(validationErrors, ValidationError{Field: "social_media_url", Message: "Social Media URL harus berupa string"})
    }

    if len(validationErrors) > 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Terjadi kesalahan dalam validasi",
            "errors":  validationErrors,
        })
        return
    }

    var sosmed models.SocialMedia
    if err := db.Where("id = ?", sosmedID).First(&sosmed).Error; err != nil {
        ctx.AbortWithError(http.StatusNotFound, errors.New("Sosial media tidak ditemukan"))
        return
    }

    updateData := make(map[string]interface{})

    if name, ok := requestData["name"].(string); ok {
        updateData["name"] = name
    }

    if SocialMediaUrl, ok := requestData["social_media_url"].(string); ok {
        updateData["social_media_url"] = SocialMediaUrl
    }

    if err := db.Model(&sosmed).Updates(updateData).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    sosmedDataJSON := map[string]interface{}{
        "id":               sosmed.Id,
        "name":             updateData["name"],
        "social_media_url": updateData["social_media_url"],
        "user_id":          sosmed.UserID,
    }

    ctx.JSON(http.StatusOK, sosmedDataJSON)
}



func DeleteSocialMediaByID(ctx *gin.Context) {
    db := database.GetDB()

    sosmedID := ctx.Param("socialmediaId")

    var sosmed models.SocialMedia
    if err := db.Where("id = ?", sosmedID).First(&sosmed).Error; err != nil {
        ctx.AbortWithError(http.StatusNotFound, errors.New("Social media tidak ditemukan"))
        return
    }

	if err := db.Delete(&sosmed).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Social media berhasil dihapus",
    })
}
