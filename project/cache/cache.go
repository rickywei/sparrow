package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/conf"
	"github.com/rickywei/sparrow/project/logger"
)

var (
	// redis cache
	RC redis.UniversalClient

	// locale cache
	LC *bigcache.BigCache
)

func init() {
	var err error
	if viper.IsSet("redis.cluster") {
		RC = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: conf.Strings("redis.cluster.addrs"),
		})
	} else if viper.IsSet("redis.sentinel") {
		RC = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:      conf.Strings("redis.sentinel.addrs"),
			MasterName: conf.String("redis.sentinel.addrs"),
		})
	} else {
		RC = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: conf.Strings("redis.client.addr"),
		})
	}
	if _, err = RC.Ping(context.Background()).Result(); err != nil {
		logger.L().Fatal("ping redis failed", zap.Error(err))
	}

	if LC, err = bigcache.New(context.Background(), bigcache.DefaultConfig(time.Minute)); err != nil {
		logger.L().Fatal("init localcache failed", zap.Error(err))
	}
}
