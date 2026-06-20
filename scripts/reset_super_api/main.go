package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/conf"

	_ "github.com/lib/pq"
)

var configFile = flag.String("f", "rpc/etc/core.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	err := conf.Load(*configFile, &c, conf.UseEnv())
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	dbEnt := ent.NewClient(
		ent.Driver(config.NewNoCacheDriver(c.DatabaseConf)),
	)
	defer dbEnt.Close()

	ctx := context.Background()

	// 1. 查询所有 API
	apis, err := dbEnt.API.Query().All(ctx)
	if err != nil {
		log.Fatalf("查询 API 失败: %v", err)
	}

	// 2. 初始化 Casbin
	casbinEnforcer := c.CasbinConf.MustNewCasbinWithRedisWatcher(c.DatabaseConf.DBType, c.DatabaseConf.GetDSN(), c.RedisConf)

	// 3. 构建策略
	var policies [][]string
	for _, apiItem := range apis {
		policies = append(policies, []string{"super", apiItem.Path, apiItem.Method})
	}

	// 4. 清理旧策略并添加新策略
	oldPolicies, err := casbinEnforcer.GetFilteredPolicy(0, "super")
	if err != nil {
		log.Fatalf("获取超级管理员旧 API 策略失败: %v", err)
	}

	if len(oldPolicies) != 0 {
		_, err = casbinEnforcer.RemoveFilteredPolicy(0, "super")
		if err != nil {
			log.Fatalf("删除超级管理员旧 API 策略失败: %v", err)
		}
	}

	if len(policies) > 0 {
		added, err := casbinEnforcer.AddPolicies(policies)
		if err != nil || !added {
			log.Fatalf("添加超级管理员新 API 策略失败: %v", err)
		}
	}

	fmt.Println("🎉 超级管理员 (super) API 权限重置成功！")
	fmt.Printf("共分配了 %d 条 API 策略。\n", len(policies))
}
