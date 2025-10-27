package user

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewGetUserInfoLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResponse, err error) {
	u, err := l.svcCtx.UserRpc.GetUser(l.ctx, &userservice.UUIDRequest{Id: l.ctx.Value("userId").(string)})
	if err != nil {
		return nil, err
	}

	resp = &types.UserInfoResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.UserInfo{
			Avatar:         pointer.GetString(u.Avatar),
			RealName:       pointer.GetString(u.FullName),
			Roles:          u.RoleValues,
			UserId:         pointer.GetString(u.Id),
			Username:       pointer.GetString(u.Username),
			Desc:           pointer.GetString(u.UserDescription),
			HomePath:       pointer.GetString(u.HomePath),
			Email:          pointer.GetString(u.Email),
			RoleNames:      u.RoleNames,
			DepartmentName: pointer.GetString(u.DepartmentName),
		},
	}

	return
}
