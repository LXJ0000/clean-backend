package middleware

import (
	"net/http"
	"strings"

	"github.com/LXJ0000/clean-backend/internal/domain"
	tokenutil "github.com/LXJ0000/clean-backend/utils/token"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := tokenutil.IsAuthorized(authToken, secret)
			if authorized {
				userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, domain.ErrorResp("ExtractIDFromToken Fail", err))
					c.Abort()
					return
				}
				c.Set(domain.UserAuthorization, userID)
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, domain.ErrorResp("Not authorized", err))
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, domain.ErrorResp("Not authorized", nil))
		c.Abort()
	}
}
