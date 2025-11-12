package role

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateRoleLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建或更新角色
func NewCreateOrUpdateRoleLogic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdateRoleLogic {
	return &CreateOrUpdateRoleLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CreateOrUpdateRoleLogic) CreateOrUpdateRole(req *types.RoleInfo) (resp *types.ModifyRoleResponse, err error) {
	// 将 API RoleInfo 转换为 RPC RoleInfo
	rpcReq := ConvertApiRoleInfoToRpcRoleInfo(req)
	rpcResp, err := l.svcCtx.RoleRpc.CreateOrUpdateRole(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	info := ConvertRpcRoleInfoToApiRoleInfo(rpcResp)
	resp = &types.ModifyRoleResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: *info,
	}
	return resp, nil
}
