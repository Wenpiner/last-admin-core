package department

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/departmentservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDepartmentLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除部门
func NewDeleteDepartmentLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteDepartmentLogic {
	return &DeleteDepartmentLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeleteDepartmentLogic) DeleteDepartment(req *types.ID32Request) (resp *types.BaseResponse, err error) {
	// 调用 RPC 服务进行删除
	rpcReq := &departmentservice.ID32Request{Id: req.ID}
	rpcResp, err := l.svcCtx.DepartmentRpc.DeleteDepartment(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	resp = &types.BaseResponse{
		Code:    0,
		Message: rpcResp.Message,
	}
	return resp, nil
}
