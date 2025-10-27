package userservicelogic

import (
	"context"

	"github.com/google/uuid"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent/user"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/userutils"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户
func (l *GetUserLogic) GetUser(in *core.UUIDRequest) (*core.UserInfo, error) {
	// 解析用户ID
	userID, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 查询用户信息（包含关联数据）
	userWithEdges, err := l.svcCtx.DBEnt.User.Query().
		Where(user.IDEQ(userID)).
		WithRoles().
		WithPositions().
		WithDepartment().
		WithTotp().
		Only(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return userutils.ConvertUserToUserInfo(userWithEdges), nil
}