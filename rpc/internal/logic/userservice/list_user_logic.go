package userservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/ent/user"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/userutils"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户列表
func (l *ListUserLogic) ListUser(in *core.UserListRequest) (*core.UserListResponse, error) {
	// 构建查询条件
	var predicates []predicate.User

	// 根据用户名模糊搜索
	if in.Username != nil && *in.Username != "" {
		predicates = append(predicates, user.UsernameContains(*in.Username))
	}

	// 根据邮箱模糊搜索
	if in.Email != nil && *in.Email != "" {
		predicates = append(predicates, user.EmailContains(*in.Email))
	}

	// 根据手机号模糊搜索
	if in.Mobile != nil && *in.Mobile != "" {
		predicates = append(predicates, user.MobileContains(*in.Mobile))
	}

	// 执行分页查询（包含关联数据）
	page, err := l.svcCtx.DBEnt.User.Query().
		Where(predicates...).
		WithRoles().
		WithPositions().
		WithTotp().
		Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	resp := &core.UserListResponse{
		Page: &core.BasePageResp{
			Total:      page.PageDetails.Total,
			PageNumber: in.Page.PageNumber,
			PageSize:   in.Page.PageSize,
		},
	}

	for _, user := range page.List {
		resp.List = append(resp.List, userutils.ConvertUserToUserInfo(user))
	}

	return resp, nil
}

