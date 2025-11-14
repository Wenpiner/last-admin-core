package main

import (
	"flag"
	"fmt"

	"github.com/wenpiner/last-admin-core/rpc/internal/config"
	apiserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/apiservice"
	configurationserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/configurationservice"
	departmentserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/departmentservice"
	dictserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/dictservice"
	initserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/initservice"
	menuserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/menuservice"
	oauthproviderserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/oauthproviderservice"
	positionserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/positionservice"
	roleserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/roleservice"
	tokenserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/tokenservice"
	userserviceServer "github.com/wenpiner/last-admin-core/rpc/internal/server/userservice"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

var configFile = flag.String("f", "etc/core.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		core.RegisterApiServiceServer(grpcServer, apiserviceServer.NewApiServiceServer(ctx))
		core.RegisterDictServiceServer(grpcServer, dictserviceServer.NewDictServiceServer(ctx))
		core.RegisterRoleServiceServer(grpcServer, roleserviceServer.NewRoleServiceServer(ctx))
		core.RegisterMenuServiceServer(grpcServer, menuserviceServer.NewMenuServiceServer(ctx))
		core.RegisterDepartmentServiceServer(grpcServer, departmentserviceServer.NewDepartmentServiceServer(ctx))
		core.RegisterPositionServiceServer(grpcServer, positionserviceServer.NewPositionServiceServer(ctx))
		core.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(ctx))
		core.RegisterOauthProviderServiceServer(grpcServer, oauthproviderserviceServer.NewOauthProviderServiceServer(ctx))
		core.RegisterInitServiceServer(grpcServer, initserviceServer.NewInitServiceServer(ctx))
		core.RegisterTokenServiceServer(grpcServer, tokenserviceServer.NewTokenServiceServer(ctx))
		core.RegisterConfigurationServiceServer(grpcServer, configurationserviceServer.NewConfigurationServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
