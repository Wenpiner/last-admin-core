package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/zeromicro/go-zero/core/logx"
)

// CacheRefreshService 缓存刷新服务
// 定时从数据库加载所有配置到缓存
type CacheRefreshService struct {
	cache           *ConfigurationCache
	dbEnt           *ent.Client
	logger          logx.Logger
	ticker          *time.Ticker
	stopChan        chan struct{}
	refreshInterval time.Duration
}

// NewCacheRefreshService 创建缓存刷新服务
func NewCacheRefreshService(cache *ConfigurationCache, dbEnt *ent.Client, logger logx.Logger) *CacheRefreshService {
	return &CacheRefreshService{
		cache:           cache,
		dbEnt:           dbEnt,
		logger:          logger,
		stopChan:        make(chan struct{}),
		refreshInterval: 5 * time.Minute,
	}
}

// Start 启动缓存刷新服务
func (s *CacheRefreshService) Start(ctx context.Context) {
	// 首先进行初始化加载
	s.refresh(ctx)

	// 启动定时刷新
	go s.backgroundRefresh(ctx)
}

// Stop 停止缓存刷新服务
func (s *CacheRefreshService) Stop() {
	close(s.stopChan)
	if s.ticker != nil {
		s.ticker.Stop()
	}
}

// backgroundRefresh 后台定时刷新
func (s *CacheRefreshService) backgroundRefresh(ctx context.Context) {
	s.ticker = time.NewTicker(s.refreshInterval)
	defer s.ticker.Stop()

	for {
		select {
		case <-s.ticker.C:
			s.refresh(ctx)
		case <-s.stopChan:
			return
		}
	}
}

// refresh 从数据库刷新所有配置到缓存
func (s *CacheRefreshService) refresh(ctx context.Context) {
	// 查询所有配置
	configs, err := s.dbEnt.Configuration.Query().All(ctx)
	if err != nil {
		s.logger.Errorf("failed to refresh configuration cache: %v", err)
		return
	}

	// 构建 key -> value 映射
	configMap := make(map[string]string)
	for _, config := range configs {
		v := fmt.Sprintf("%s<>%s", config.Group, config.Value)
		configMap[config.Key] = v
	}

	// 更新缓存
	s.cache.SetAll(configMap)
	s.logger.Infof("configuration cache refreshed, total configs: %d", len(configMap))
}

// RefreshNow 立即刷新缓存（用于 Create/Update/Delete 后调用）
func (s *CacheRefreshService) RefreshNow(ctx context.Context) {
	s.refresh(ctx)
}
