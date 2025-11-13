package configurationservicelogic

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/types/core"
)

// ConvertConfigurationToConfigurationInfo 将 Configuration 实体转换为 ConfigurationInfo
func ConvertConfigurationToConfigurationInfo(config *ent.Configuration) *core.ConfigurationInfo {
	return &core.ConfigurationInfo{
		Key:         config.Key,
		Value:       config.Value,
		Name:        config.Name,
		Group:       config.Group,
		Description: pointer.ToStringPtrIfNotEmpty(config.Description),
	}
}

