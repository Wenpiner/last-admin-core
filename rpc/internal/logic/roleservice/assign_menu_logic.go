package roleservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignMenuLogic {
	return &AssignMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 为角色分配菜单
func (l *AssignMenuLogic) AssignMenu(in *core.RoleMenuRequest) (*core.BaseResponse, error) {
	err := l.svcCtx.DBEnt.Role.UpdateOneID(*in.RoleId).ClearMenus().AddMenuIDs(in.MenuIds...).Exec(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	return &core.BaseResponse{
		Message: "success",
	}, nil
}
