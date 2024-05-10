package bootstrap

import (
	"github.com/IBM/sarama"
	"github.com/LXJ0000/clean-backend/internal/event"
	"github.com/LXJ0000/clean-backend/utils/cache"
	logutil "github.com/LXJ0000/clean-backend/utils/log"
	"github.com/LXJ0000/clean-backend/utils/orm"
	snowflakeutil "github.com/LXJ0000/clean-backend/utils/snowflake"
)

type Application struct {
	Env *Env

	Orm        orm.Database
	RedisCache cache.RedisCache

	SaramaClient sarama.Client
	Producer     event.Producer
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Orm = NewOrmDatabase(app.Env)
	app.RedisCache = NewRedisCache(app.Env)

	app.SaramaClient = NewSaramaClient(app.Env)
	app.Producer = NewProducer(app.Env)
	
	logutil.Init(app.Env.AppEnv)
	snowflakeutil.Init(app.Env.SnowflakeStartTime, app.Env.SnowflakeMachineID)

	return *app
}
