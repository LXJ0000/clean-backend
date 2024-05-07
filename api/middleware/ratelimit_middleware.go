package middleware

import (
	"github.com/LXJ0000/go-lib/gin-plugin/ratelimit"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	return ratelimit.New(
		ratelimit.DefaultConfig(),
	)
}
