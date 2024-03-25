package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "4hsh4du"

func GenerateToken(id uint, email string) (token string, err error) {

	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = parseToken.SignedString([]byte(secretKey))

	return
}

func VerifyToken(ctx *gin.Context) (interface{}, error) {
	
	jwt_token := ctx.Request.Header.Get("Authorization")

	bearer := strings.HasPrefix(jwt_token, "Bearer")
	if !bearer {
		return nil, errors.New("Bearer token not found")
	}

	tokenString := strings.TrimSpace(jwt_token[7:])

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Metode penanda tangan tidak sesuai yang diharapkan")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.New("Invalid token")
	}
	

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		return nil, errors.New("Failed to parse claims")
	}

	return token.Claims.(jwt.MapClaims), nil
}

