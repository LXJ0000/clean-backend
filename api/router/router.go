package route

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/clean-backend/bootstrap"
	"github.com/LXJ0000/clean-backend/internal/event"
	"github.com/LXJ0000/clean-backend/utils/cache"
	"github.com/LXJ0000/clean-backend/utils/orm"
	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db orm.Database, redisCache cache.RedisCache, gin *gin.Engine, saramaClient sarama.Client, producer event.Producer) {

	// publicRouter := gin.Group("/api")
	// All Public APIs

	// protectedRouter := gin.Group("/api")
	// // Middleware to verify AccessToken
	// protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs

	// User
	NewUserRouter(env, timeout, db, gin)
	NewPostRouter(env, timeout, db, gin, redisCache, saramaClient, producer)
}
