package controllers

import (
	"errors"
	"net/http"
	"strings"

    "main.go/database"
	"main.go/helpers"
	"main.go/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}


func CreateUser(ctx *gin.Context) {
    db := database.GetDB()
    user := models.User{}

    if err := ctx.ShouldBindJSON(&user)
	err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Gagal memvalidasi data",
            "error":   err.Error(),
        })
        return
    }

	var validationErrors []ValidationError

    if user.Email == "" {
        validationErrors = append(validationErrors, ValidationError{Field: "email", Message: "Email harus diisi"})
    } else if strings.Count(user.Email, "@") != 1 {
        validationErrors = append(validationErrors, ValidationError{Field: "email", Message: "Format Email harus sesuai"})
    }

	if user.Password == "" {
	validationErrors = append(validationErrors, ValidationError{Field: "password", Message: "Email harus diisi"})
	}

	if user.Age < 8 {
	validationErrors = append(validationErrors, ValidationError{Field: "age", Message: "Umur minimal 8 Tahun"})
	}

	if !strings.HasPrefix(user.ProfileImageUrl, "http://") && !strings.HasPrefix(user.ProfileImageUrl, "https://") {
	validationErrors = append(validationErrors, ValidationError{Field: "profile_image_url", Message: "Format profile_image_url harus diawali dengan http:// atau https://"})
	}

	if len(validationErrors) > 0 {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": "Terjadi kesalahan dalam validasi",
		"errors":  validationErrors,
	})
	return
	}

	if !strings.HasPrefix(user.ProfileImageUrl, "http://") && !strings.HasPrefix(user.ProfileImageUrl, "https://") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Format profile_image_url harus diawali dengan http:// atau https://",
		})
		return
	}
	
    _, err := govalidator.ValidateStruct(&user)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Validasi gagal",
            "error":   err.Error(),
        })
        return
    }

    if err := db.Create(&user).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Gagal membuat pengguna",
            "error":   err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{
        "id":               user.Id,
        "email":            user.Email,
        "username":         user.Username,
        "age":              user.Age,
        "profile_image_url": user.ProfileImageUrl,
    })
}


func LoginUser(ctx *gin.Context){
	
	db := database.GetDB()
	user := models.User{}


	if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Gagal memvalidasi data",
            "error":   err.Error(),
        })
        return
    }

	pwd := user.Password

	err := db.Where("email = ?", user.Email).Take(&user).Error
	if err != nil{
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !helpers.PasswordValid(user.Password, pwd){
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid password"))
		return
	}

	token, err := helpers.GenerateToken(user.Id, user.Email)
	if err != nil{
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	
    jwt_token := "Bearer " + token



	ctx.JSON(http.StatusCreated, gin.H{
		"token": jwt_token,
	})

}

func ValidateTes (ctx *gin.Context) {
    ctx.JSON(http.StatusOK, gin.H{
        "msg": "AUTENTIKASI BERHASIL",
    })
}

func UpdateUser(ctx *gin.Context) {
    db := database.GetDB()

    userData := ctx.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))

    updateData := make(map[string]interface{})

    var requestData map[string]interface{}
    if err := ctx.ShouldBindJSON(&requestData); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Gagal memvalidasi data",
            "error":   err.Error(),
        })
        return
    }

    var validationErrors []ValidationError

    if email, ok := requestData["email"].(string); ok {
        if email == "" {
            validationErrors = append(validationErrors, ValidationError{Field: "email", Message: "Email harus diisi"})
        } else if !strings.Contains(email, "@") || strings.Count(email, "@") > 1 {
            validationErrors = append(validationErrors, ValidationError{Field: "email", Message: "Email harus menggunakan satu '@' dan tidak lebih dari satu"})
        }
        updateData["email"] = email
    }

    if username, ok := requestData["username"].(string); ok {
        if username == "" {
            validationErrors = append(validationErrors, ValidationError{Field: "username", Message: "Username harus diisi"})
        }
        updateData["username"] = username
    }

    if ageFloat, ok := requestData["age"].(float64); ok {
        age := int(ageFloat)
        if age <= 8 {
            validationErrors = append(validationErrors, ValidationError{Field: "age", Message: "Umur harus lebih besar dari 8"})
        }
        updateData["age"] = age
    }

    if profileImageURL, ok := requestData["profile_image_url"].(string); ok {
        if profileImageURL != "" && !strings.HasPrefix(profileImageURL, "http://") && !strings.HasPrefix(profileImageURL, "https://") {
            validationErrors = append(validationErrors, ValidationError{Field: "profile_image_url", Message: "Profile image URL harus menggunakan http:// atau https://"})
        }
        updateData["profile_image_url"] = profileImageURL
    }

    if len(validationErrors) > 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Terjadi kesalahan dalam validasi",
            "errors":  validationErrors,
        })
        return
    }

    if err := db.Model(&models.User{}).Where("id = ?", userID).Updates(updateData).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    userDataJSON := map[string]interface{}{
        "id":               userID,
        "email":            updateData["email"],
        "username":         updateData["username"],
        "age":              updateData["age"],
        "profile_image_url": updateData["profile_image_url"],
    }

    ctx.JSON(http.StatusOK, userDataJSON)
}

func DeleteUser(ctx *gin.Context) {
    userData := ctx.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))

    db := database.GetDB()
    if err := db.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
        ctx.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Data pengguna berhasil dihapus",
    })
}