package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (hdl *Handler) JWTAuth(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	if authHeader == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	type Claims struct {
		UserID int
		jwt.RegisteredClaims
	}

	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(hdl.cfg.JWT_SECRET), nil
	})

	if err != nil || !token.Valid {
		fmt.Println(err, token.Valid)
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	context.Set("ID", claims.UserID)
}

func (hdl *Handler) StatusLogger(context *gin.Context) {

}
