package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vonrimez/TaskAPI/internal/auth"
)

func (hdl *Handler) JWTAuth(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	if authHeader == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, "no token")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := auth.IsValidJWT(tokenString, hdl.userService.JWTSecret)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, "internal error")
		return
	} else {
		if claims == nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
	}

	context.Set("ID", claims.UserID)
}

func (hdl *Handler) StatusLogger(context *gin.Context) {
	context.Next()
	fmt.Println("[ERR_LOGER]", context.Err())
}
