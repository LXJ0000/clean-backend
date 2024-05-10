package route

import (
	"log"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/clean-backend/api/handler"
	"github.com/LXJ0000/clean-backend/api/middleware"
	"github.com/LXJ0000/clean-backend/bootstrap"
	"github.com/LXJ0000/clean-backend/internal/event"
	"github.com/LXJ0000/clean-backend/internal/repository"
	"github.com/LXJ0000/clean-backend/internal/service"
	"github.com/LXJ0000/clean-backend/utils/cache"
	"github.com/LXJ0000/clean-backend/utils/orm"
	"github.com/gin-gonic/gin"
)

func NewPostRouter(env *bootstrap.Env, timeout time.Duration, db orm.Database, r *gin.Engine, redisCache cache.RedisCache, saramaClient sarama.Client, producer event.Producer) {
	postRepo := repository.NewPostRepository(db, redisCache)
	interRepo := repository.NewInteractionRepository(db, redisCache)
	interactionService := service.NewInteractionService(interRepo, timeout)
	postHandler := &handler.PostHandler{
		PostService: service.NewPostService(postRepo, timeout, env, interactionService, producer),
	}

	consumer := event.NewBatchSyncReadEventConsumer(saramaClient, interRepo)
	if err := consumer.Start(); err != nil {
		slog.Error("OMG！消费者启动失败")
		log.Fatal(err)
	}

	post := r.Group("/post")
	post.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))

	post.GET("/detail", postHandler.Detail)
	post.POST("/create", postHandler.Create)
	post.POST("/publish", postHandler.Publish)
	post.GET("/list/reader", postHandler.ReaderList)
	post.GET("/list/writer", postHandler.WriterList)
}
