package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/role"
	"github.com/wenpiner/last-admin-core/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/conf"

	_ "github.com/lib/pq"
)

var configFile = flag.String("f", "rpc/etc/core.yaml", "the config file")
var roleCode = flag.String("role", "super", "the role code to reset API permissions for")

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

	// 0. 查询角色是否存在
	roleEntity, err := dbEnt.Role.Query().Where(role.RoleCodeEQ(*roleCode)).Only(ctx)
	if err != nil {
		log.Fatalf("查询角色编码 '%s' 失败或该角色不存在: %v", *roleCode, err)
	}

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
		policies = append(policies, []string{roleEntity.RoleCode, apiItem.Path, apiItem.Method})
	}

	// 4. 清理旧策略并添加新策略
	oldPolicies, err := casbinEnforcer.GetFilteredPolicy(0, roleEntity.RoleCode)
	if err != nil {
		log.Fatalf("获取角色 '%s' 旧 API 策略失败: %v", roleEntity.RoleCode, err)
	}

	if len(oldPolicies) != 0 {
		_, err = casbinEnforcer.RemoveFilteredPolicy(0, roleEntity.RoleCode)
		if err != nil {
			log.Fatalf("删除角色 '%s' 旧 API 策略失败: %v", roleEntity.RoleCode, err)
		}
	}

	if len(policies) > 0 {
		added, err := casbinEnforcer.AddPolicies(policies)
		if err != nil || !added {
			log.Fatalf("添加角色 '%s' 新 API 策略失败: %v", roleEntity.RoleCode, err)
		}
	}

	fmt.Printf("🎉 角色 '%s' (%s) API 权限重置成功！\n", roleEntity.RoleName, roleEntity.RoleCode)
	fmt.Printf("共分配了 %d 条 API 策略。\n", len(policies))
}
