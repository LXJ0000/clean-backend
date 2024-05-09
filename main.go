package main

import (
	"time"

	"github.com/LXJ0000/clean-backend/api/middleware"
	route "github.com/LXJ0000/clean-backend/api/router"
	"github.com/LXJ0000/clean-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Orm
	redisCache := app.RedisCache

	timeout := time.Duration(env.ContextTimeout) * time.Hour // 接口超时时间

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())
	server.Use(middleware.RateLimitMiddleware())
	route.Setup(env, timeout, db, redisCache, server)

	_ = server.Run(env.ServerAddress)

}
