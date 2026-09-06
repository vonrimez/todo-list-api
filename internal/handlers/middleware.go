package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/vonrimez/TaskAPI/internal/auth"
)

func (hdl *Handler) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithJSON(w, &Response{
				statusCode: http.StatusUnauthorized,
				body:       H{"error": "no token"},
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := auth.IsValidJWT(tokenString, hdl.userService.JWTSecret)
		if err != nil {
			respondWithJSON(w, &Response{
				statusCode: http.StatusInternalServerError,
				body:       H{"error": "internal error"},
			})
			return
		} else {
			if claims == nil {
				respondWithJSON(w, &Response{
					statusCode: http.StatusUnauthorized,
					body:       H{"error": "invalid token"},
				})
				return
			}
		}

		setValue(r, UserID, claims.UserID)
		log.Println(r.Context().Value(UserID))
		next.ServeHTTP(w, r)
	})
}
