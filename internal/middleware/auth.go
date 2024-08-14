package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samdyra/go-geo/internal/utils/auth"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
            c.Abort()
            return
        }

        userID, username, err := auth.ValidateToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, errors.NewAPIError(errors.ErrUnauthorized))
            c.Abort()
            return
        }

        c.Set("user_id", userID)
        c.Set("username", username)
        c.Next()
    }
}