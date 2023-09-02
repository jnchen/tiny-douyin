package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func RateLimit(limit rate.Limit, burst int) gin.HandlerFunc {
	interval := 2 * time.Duration(float64(time.Second)/float64(limit))
	limiter := rate.NewLimiter(limit, burst)
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), interval)
		defer cancel()

		if err := limiter.Wait(ctx); err != nil {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

func QPSLimit(qps int) gin.HandlerFunc {
	return RateLimit(rate.Limit(qps), qps)
}
