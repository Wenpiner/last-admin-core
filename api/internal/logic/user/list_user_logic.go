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

type ListUserLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户列表
func NewListUserLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListUserLogic) ListUser(req *types.UserListRequest) (resp *types.UserListResponse, err error) {
	// 构建 RPC 请求
	rpcReq := &userservice.UserListRequest{
		Page: &userservice.BasePageRequest{
			PageNumber: req.Page.CurrentPage,
			PageSize:   req.Page.PageSize,
		},
		Username: pointer.ToStringPtrIfNotEmpty(req.Username),
		Mobile:   pointer.ToStringPtrIfNotEmpty(req.Mobile),
		Email:    pointer.ToStringPtrIfNotEmpty(req.Email),
	}

	// 调用 RPC 服务获取用户列表
	rpcResp, err := l.svcCtx.UserRpc.ListUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	resp = ConvertRpcUserListResponseToApiUserListResponse(rpcResp)
	return resp, nil
}
