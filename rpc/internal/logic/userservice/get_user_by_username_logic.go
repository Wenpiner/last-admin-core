package userservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/user"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/userutils"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByUsernameLogic {
	return &GetUserByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户 - 用户用户名
func (l *GetUserByUsernameLogic) GetUserByUsername(in *core.StringRequest) (*core.UserInfo, error) {
	// 查询用户信息（包含关联数据）
	userWithEdges, err := l.svcCtx.DBEnt.User.Query().
		WithRoles().
		WithPositions().
		WithDepartment().
		WithTotp().
		Where(user.UsernameEQ(in.Value)).
		Only(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return userutils.ConvertUserToUserInfo(userWithEdges), nil
}
