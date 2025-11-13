package configuration

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateConfigurationLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增/更新配置
func NewCreateOrUpdateConfigurationLogic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdateConfigurationLogic {
	return &CreateOrUpdateConfigurationLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CreateOrUpdateConfigurationLogic) CreateOrUpdateConfiguration(req *types.ConfigurationInfo) (resp *types.ConfigurationInfo, err error) {
	// 将 API ConfigurationInfo 转换为 RPC ConfigurationInfo
	rpcReq := ConvertApiConfigurationInfoToRpcConfigurationInfo(req)

	// 调用 RPC 服务创建或更新配置
	rpcResp, err := l.svcCtx.ConfigurationRpc.CreateOrUpdateConfiguration(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	resp = ConvertRpcConfigurationInfoToApiConfigurationInfo(rpcResp)
	return resp, nil
}
