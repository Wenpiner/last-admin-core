package roleservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/role"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取角色菜单
func (l *GetMenuLogic) GetMenu(in *core.ID32Request) (*core.RoleMenuListResponse, error) {
	menus, err := l.svcCtx.DBEnt.Role.Query().Where(role.IDEQ(in.Id)).WithMenus().QueryMenus().IDs(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	return &core.RoleMenuListResponse{
		List: menus,
	}, nil
}
