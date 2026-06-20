package role

import (
	"context"
	"os"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type ResetSuperAdminApiLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置超级管理员API权限
func NewResetSuperAdminApiLogic(r *http.Request, svcCtx *svc.ServiceContext) *ResetSuperAdminApiLogic {
	return &ResetSuperAdminApiLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ResetSuperAdminApiLogic) ResetSuperAdminApi(req *types.ID32Request) (resp *types.BaseResponse, err error) {
	if os.Getenv("ALLOW_RESET_SUPER_API") != "true" {
		return nil, errorx.NewApiError(errorx.CodeForbidden, "功能未开启，请在后端配置 ALLOW_RESET_SUPER_API=true 环境变量")
	}

	rpcReq := &core.ID32Request{
		Id: req.ID,
	}
	rpcResp, err := l.svcCtx.RoleRpc.ResetSuperApi(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: rpcResp.Message,
	}

	return
}
