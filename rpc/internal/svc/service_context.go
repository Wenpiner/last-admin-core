package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	DBEnt  *ent.Client
	Redis redis.UniversalClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	dbEnt := ent.NewClient(
		ent.Log(logx.Error),
		ent.Driver(config.NewNoCacheDriver(c.DatabaseConf)),
	)
	return &ServiceContext{
		Config: c,
		DBEnt:  dbEnt,
		Redis:  c.RedisConf.NewMustUniversalRedis(),
	}
}
