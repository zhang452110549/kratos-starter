package data

import (
	"context"

	"kratos-starter/internal/conf"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *conf.Data) redis.UniversalClient {
	rdCfg := cfg.GetRedis()
	if rdCfg == nil {
		panic("redis config is nil")
	}
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            rdCfg.GetAddresses(),
		DB:               int(rdCfg.GetDb()),
		Username:         rdCfg.GetUserName(),
		Password:         rdCfg.GetPassword(),
		SentinelUsername: rdCfg.GetSentinelUserName(),
		SentinelPassword: rdCfg.GetSentinelPassword(),
		ReadTimeout:      rdCfg.GetReadTimeout().AsDuration(),
		WriteTimeout:     rdCfg.GetWriteTimeout().AsDuration(),
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	return rdb
}
