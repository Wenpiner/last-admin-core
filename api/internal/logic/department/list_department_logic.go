package department

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/departmentservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDepartmentLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取部门列表
func NewListDepartmentLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListDepartmentLogic {
	return &ListDepartmentLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListDepartmentLogic) ListDepartment(req *types.DepartmentListRequest) (resp *types.DepartmentListResponse, err error) {
	// 构建 RPC 请求
	rpcReq := &departmentservice.DepartmentListRequest{
		Page: &departmentservice.BasePageRequest{
			PageNumber: req.Page.CurrentPage,
			PageSize:   req.Page.PageSize,
		},
		DeptName: pointer.ToStringPtr(req.DeptName),
		DeptCode: pointer.ToStringPtr(req.DeptCode),
	}

	// 调用 RPC 服务获取部门列表
	rpcResp, err := l.svcCtx.DepartmentRpc.ListDepartment(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	resp = ConvertRpcDepartmentListResponseToApiDepartmentListResponse(rpcResp)
	return resp, nil
}
