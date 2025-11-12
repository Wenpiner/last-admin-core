package svc

import (
	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	DBEnt  *ent.Client
	Redis  redis.UniversalClient
	Casbin *casbin.Enforcer
}

func NewServiceContext(c config.Config) *ServiceContext {
	dbEnt := ent.NewClient(
		ent.Log(logx.Error),
		ent.Driver(config.NewNoCacheDriver(c.DatabaseConf)),
	)

	casbin, err := c.CasbinConf.NewCasbin(c.DatabaseConf.DBType, c.DatabaseConf.GetDSN())
	if err != nil {
		logx.Errorw("初始化 Casbin 失败", logx.Field("detail", err.Error()))
	}

	return &ServiceContext{
		Config: c,
		DBEnt:  dbEnt,
		Redis:  c.RedisConf.NewMustUniversalRedis(),
		Casbin: casbin,
	}
}
