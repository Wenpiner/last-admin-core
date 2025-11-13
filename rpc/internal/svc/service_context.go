package svc

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/cache"
	"github.com/wenpiner/last-admin-core/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config              config.Config
	DBEnt               *ent.Client
	Redis               redis.UniversalClient
	Casbin              *casbin.Enforcer
	ConfigurationCache  *cache.ConfigurationCache
	CacheRefreshService *cache.CacheRefreshService
	ConfigValidator     *cache.ConfigValidator
}

func NewServiceContext(c config.Config) *ServiceContext {
	dbEnt := ent.NewClient(
		ent.Log(logx.Error),
		ent.Driver(config.NewNoCacheDriver(c.DatabaseConf)),
	)

	casbin := c.CasbinConf.MustNewCasbinWithRedisWatcher(c.DatabaseConf.DBType, c.DatabaseConf.GetDSN(), c.RedisConf)

	// 初始化配置缓存
	configCache := cache.NewConfigurationCache()
	cacheRefreshService := cache.NewCacheRefreshService(configCache, dbEnt, logx.WithContext(context.Background()))

	// 启动缓存刷新服务
	cacheRefreshService.Start(context.Background())

	// 初始化配置验证器
	configValidator, err := cache.NewConfigValidator(configCache)
	if err != nil {
		logx.Errorw("初始化 ConfigValidator 失败", logx.Field("detail", err.Error()))
	}

	return &ServiceContext{
		Config:              c,
		DBEnt:               dbEnt,
		Redis:               c.RedisConf.NewMustUniversalRedis(),
		Casbin:              casbin,
		ConfigurationCache:  configCache,
		CacheRefreshService: cacheRefreshService,
		ConfigValidator:     configValidator,
	}
}
