package config

import (
	"github.com/wenpiner/last-admin-common/config"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-common/plugins/casbin"
	"github.com/wenpiner/last-admin-common/utils/captcha"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	CaptchaConf        captcha.CaptchaConfig // 验证码配置
	ProjectConf        ProjectConfig         // 项目配置
	I18nConf           last_i18n.Config      // 国际化配置
	CoreRpc            zrpc.RpcClientConf
	RedisConf          config.RedisConfig
	CasbinConf         casbin.CasbinConf     `json:",optional"`
	CasbinDatabaseConf config.DatabaseConfig `json:",optional"`
}

type ProjectConfig struct {
	RegisterRoleValue string `json:",default=admin"` // 注册用户默认角色
}
