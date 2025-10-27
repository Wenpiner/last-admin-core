package config

import (
	"context"
	"database/sql"
	"time"

	esql "entgo.io/ent/dialect/sql"
	"github.com/wenpiner/last-admin-common/config"
	"github.com/wenpiner/last-admin-common/plugins/casbin"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DatabaseConf     config.DatabaseConfig `json:",optional"`               // 数据库配置
	OAuthStateSecret string                `json:",env=OAUTH_STATE_SECRET"` // OAuth 状态密钥
	RedisConf        config.RedisConfig
	CasbinConf       casbin.CasbinConf   `json:",optional"`
}

// NewNoCacheDriver returns a new driver with no cache.
func NewNoCacheDriver(c config.DatabaseConfig) *esql.Driver {
	db, err := sql.Open(c.DBType, c.GetDSN())
	logx.Must(err)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.PingContext(ctx)
	logx.Must(err)

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)

	return esql.OpenDB(c.DBType, db)
}