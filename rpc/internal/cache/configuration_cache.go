package cache

import (
	"sync"
)

// ConfigurationCache 配置缓存管理器
// 使用 Map 存储 key -> value 的映射
type ConfigurationCache struct {
	mu    sync.RWMutex
	cache map[string]string
}

// NewConfigurationCache 创建新的配置缓存
func NewConfigurationCache() *ConfigurationCache {
	return &ConfigurationCache{
		cache: make(map[string]string),
	}
}

// Get 获取配置值
func (c *ConfigurationCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.cache[key]
	return val, ok
}

// Set 设置配置值
func (c *ConfigurationCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
}

// Delete 删除配置
func (c *ConfigurationCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

// SetAll 批量设置配置（用于初始化或刷新）
func (c *ConfigurationCache) SetAll(configs map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = configs
}

// GetAll 获取所有配置
func (c *ConfigurationCache) GetAll() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// 返回副本，防止外部修改
	result := make(map[string]string)
	for k, v := range c.cache {
		result[k] = v
	}
	return result
}

// Clear 清空所有缓存
func (c *ConfigurationCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]string)
}

// Size 获取缓存大小
func (c *ConfigurationCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

