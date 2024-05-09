package route

import (
	"time"

	"github.com/LXJ0000/clean-backend/api/handler"
	"github.com/LXJ0000/clean-backend/api/middleware"
	"github.com/LXJ0000/clean-backend/bootstrap"
	"github.com/LXJ0000/clean-backend/internal/repository"
	"github.com/LXJ0000/clean-backend/internal/service"
	"github.com/LXJ0000/clean-backend/utils/orm"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db orm.Database, r *gin.Engine) {
	userRepo := repository.NewUserRepository(db)
	userHandler := &handler.UserHandler{
		UserService: service.NewUserService(userRepo, env, timeout),
	}
	user := r.Group("/user")
	user.POST("/login", userHandler.Login)
	user.POST("/signup", userHandler.Signup)

	user.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	user.GET("/detail", userHandler.Detail)
}
