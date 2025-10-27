package roleservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/role"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleByValueLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleByValueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleByValueLogic {
	return &GetRoleByValueLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过值获取角色
func (l *GetRoleByValueLogic) GetRoleByValue(in *core.StringRequest) (*core.RoleInfo, error) {
	role, err := l.svcCtx.DBEnt.Role.Query().Where(role.RoleCodeEQ(in.Value)).Only(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return ConvertRoleToRoleInfo(role), nil
}
