package middlewares

import (
	"intelliquiz/src/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiters = make(map[string]*rate.Limiter)

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if _, exists := limiters[ip]; !exists {
			limiters[ip] = rate.NewLimiter(1, 5)
		}
		limiter := limiters[ip]

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, types.TooManyRequestsErrorResponseStruct{
				StatusCode: http.StatusTooManyRequests,
				Success:    false,
				Message:    "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
