package bootstrap

import (
	"github.com/LXJ0000/clean-backend/utils/cache"
	logutil "github.com/LXJ0000/clean-backend/utils/log"
	"github.com/LXJ0000/clean-backend/utils/orm"
	snowflakeutil "github.com/LXJ0000/clean-backend/utils/snowflake"
)

type Application struct {
	Env *Env

	Orm        orm.Database
	RedisCache cache.Cache
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Orm = NewOrmDatabase(app.Env)
	app.RedisCache = NewRedisCache(app.Env)
	logutil.Init(app.Env.AppEnv)
	snowflakeutil.Init(app.Env.SnowflakeStartTime, app.Env.SnowflakeMachineID)

	return *app
}
