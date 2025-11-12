package role

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/roleservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRoleLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色列表
func NewListRoleLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListRoleLogic) ListRole(req *types.RoleListRequest) (resp *types.RoleListResponse, err error) {
	// 构建 RPC 请求
	rpcReq := &roleservice.RoleListRequest{
		Page: &roleservice.BasePageRequest{
			PageNumber: req.Page.CurrentPage,
			PageSize:   req.Page.PageSize,
		},
		RoleName: pointer.ToStringPtrIfNotEmpty(req.RoleName),
		RoleCode: pointer.ToStringPtrIfNotEmpty(req.RoleCode),
	}

	// 调用 RPC 服务获取角色列表
	rpcResp, err := l.svcCtx.RoleRpc.ListRole(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	resp = convertRpcRoleListResponseToApiRoleListResponse(rpcResp)
	return resp, nil
}

// convertRpcRoleListResponseToApiRoleListResponse 将 RPC RoleListResponse 转换为 API RoleListResponse
func convertRpcRoleListResponseToApiRoleListResponse(rpcResp *roleservice.RoleListResponse) *types.RoleListResponse {
	if rpcResp == nil {
		return nil
	}

	apiList := make([]types.RoleInfo, 0, len(rpcResp.List))
	for _, role := range rpcResp.List {
		apiList = append(apiList, *ConvertRpcRoleInfoToApiRoleInfo(role))
	}

	return &types.RoleListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.RoleListInfo{
			BaseListInfo: types.BaseListInfo{
				Total: rpcResp.Page.Total,
			},
			List: apiList,
		},
	}
}
