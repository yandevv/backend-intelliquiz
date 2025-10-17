package middlewares

import (
	"intelliquiz/src/auth"
	"intelliquiz/src/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func bearerFromHeader(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	if after, ok := strings.CutPrefix(h, "Bearer "); ok {
		return after
	}

	return ""
}

func JWTTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := bearerFromHeader(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ForbiddenErrorResponseStruct{
				StatusCode: http.StatusForbidden,
				Success:    false,
				Message:    "Forbidden",
			})
			return
		}

		claims, err := auth.ParseAccess(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ForbiddenErrorResponseStruct{
				StatusCode: http.StatusUnauthorized,
				Success:    false,
				Message:    "Forbidden",
			})
			return
		}

		c.Set("userID", claims.Subject)
		c.Next()
	}
}
