package svc

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-common/utils/captcha"
	"github.com/wenpiner/last-admin-common/validator"
	"github.com/wenpiner/last-admin-core/api/internal/config"
	"github.com/wenpiner/last-admin-core/api/internal/i18n"
	"github.com/wenpiner/last-admin-core/api/internal/middleware"
	"github.com/wenpiner/last-admin-core/rpc/client/initservice"
	"github.com/wenpiner/last-admin-core/rpc/client/menuservice"
	"github.com/wenpiner/last-admin-core/rpc/client/oauthproviderservice"
	"github.com/wenpiner/last-admin-core/rpc/client/tokenservice"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware

	CaptchaService *captcha.Service

	Trans *last_i18n.Translator

	UserRpc  userservice.UserService
	TokenRpc tokenservice.TokenService
	OauthRpc oauthproviderservice.OauthProviderService
	InitRpc  initservice.InitService
	MenuRpc  menuservice.MenuService

	validator *validator.Validator

	Redis  *redis.Client
	Casbin *casbin.Enforcer
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化验证码服务
	captchaService, err := captcha.NewService(&c.CaptchaConf)
	if err != nil {
		logx.Errorf("Failed to initialize captcha service: %v", err)
		// 如果验证码服务初始化失败，使用默认配置
		defaultConfig := captcha.DefaultConfig()
		captchaService, _ = captcha.NewService(defaultConfig)
	}

	// 初始化国际化服务
	trans := last_i18n.NewTranslator(c.I18nConf, i18n.LocaleFS)

	// 初始化翻译器
	en := en.New()
	zh := zh.New()
	uni := ut.New(en, zh)
	validator := validator.NewValidator(trans, uni)
	httpx.SetValidator(validator)

	// 初始化Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     strings.Split(c.RedisConf.Host, ",")[0],
		Password: c.RedisConf.Password,
		DB:       c.RedisConf.DB,
		PoolSize: c.RedisConf.PoolSize,
	})

	// 初始化Casbin
	casbin := c.CasbinConf.MustNewCasbinWithRedisWatcher(c.CasbinDatabaseConf.DBType, c.CasbinDatabaseConf.GetDSN(), c.RedisConf)

	// 初始化用户服务
	coreRpc := zrpc.NewClientIfEnable(c.CoreRpc)
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(trans, casbin).Handle,
		CaptchaService: captchaService,
		Trans:          trans,
		UserRpc:        userservice.NewUserService(coreRpc),
		TokenRpc:       tokenservice.NewTokenService(coreRpc),
		OauthRpc:       oauthproviderservice.NewOauthProviderService(coreRpc),
		InitRpc:        initservice.NewInitService(coreRpc),
		MenuRpc:        menuservice.NewMenuService(coreRpc),
		Redis:          redisClient,
		Casbin:         casbin,
	}
}
