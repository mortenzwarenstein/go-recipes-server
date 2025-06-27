package middleware

import (
	"github.com/gin-gonic/gin"
	"go-recipes-server/internal/dto"
	"go-recipes-server/internal/util"
	"net/http"
	"strings"
)

// JWTMiddleware validates the JWT and injects claims into context
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Error:   "Invalid cookie",
				Message: "Invalid cookie, please login first",
			})
			return
		}

		tokenString := cookie.Value
		tokenString = strings.TrimSpace(tokenString)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Error:   "Invalid header",
				Message: "Invalid header, access token is missing",
			})
			return
		}

		claims, err := util.VerifyAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Error:   err.Error(),
				Message: "Invalid token",
			})
			return
		}

		// Inject claims into Gin context for later use
		c.Set("email", claims.Email)
		c.Set("userID", claims.Subject)
		c.Next()
	}
}
