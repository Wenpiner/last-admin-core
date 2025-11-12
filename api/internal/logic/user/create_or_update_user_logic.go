package user

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateUserLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建或更新用户
func NewCreateOrUpdateUserLogic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdateUserLogic {
	return &CreateOrUpdateUserLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CreateOrUpdateUserLogic) CreateOrUpdateUser(req *types.UserInfo) (resp *types.ModifyUserResponse, err error) {
	// 将 API UserInfo 转换为 RPC UserInfo
	rpcReq := ConvertApiUserInfoToRpcUserInfo(req)

	var rpcResp *userservice.UserInfo
	// 根据是否有 ID 来判断是创建还是更新
	if req.UserId != "" {
		// 更新用户
		rpcResp, err = l.svcCtx.UserRpc.UpdateUser(l.ctx, rpcReq)
	} else {
		// 创建用户
		rpcResp, err = l.svcCtx.UserRpc.CreateUser(l.ctx, rpcReq)
	}

	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	info := ConvertRpcUserInfoToApiUserInfo(rpcResp)
	resp = &types.ModifyUserResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: *info,
	}
	return resp, nil
}
