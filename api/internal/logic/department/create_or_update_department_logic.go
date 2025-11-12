package department

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateDepartmentLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建或更新部门
func NewCreateOrUpdateDepartmentLogic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdateDepartmentLogic {
	return &CreateOrUpdateDepartmentLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CreateOrUpdateDepartmentLogic) CreateOrUpdateDepartment(req *types.DepartmentInfo) (resp *types.ModifyDepartmentResponse, err error) {
	// 调用 RPC 服务进行创建或更新
	rpcReq := ConvertApiDepartmentInfoToRpcDepartmentInfo(req)
	rpcResp, err := l.svcCtx.DepartmentRpc.CreateOrUpdateDepartment(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	department := ConvertRpcDepartmentInfoToApiDepartmentInfo(rpcResp)
	resp = &types.ModifyDepartmentResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: *department,
	}
	return resp, nil
}
