package configuration

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConfigurationLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取配置列表
func NewListConfigurationLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListConfigurationLogic {
	return &ListConfigurationLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListConfigurationLogic) ListConfiguration(req *types.ConfigurationListRequest) (resp *types.ConfigurationListResponse, err error) {
	// 将 API ConfigurationListRequest 转换为 RPC ConfigurationListRequest
	rpcReq := ConvertApiConfigurationListRequestToRpcConfigurationListRequest(req)

	// 调用 RPC 服务获取配置列表
	rpcResp, err := l.svcCtx.ConfigurationRpc.ListConfiguration(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	resp = ConvertRpcConfigurationListResponseToApiConfigurationListResponse(rpcResp)
	return resp, nil
}
