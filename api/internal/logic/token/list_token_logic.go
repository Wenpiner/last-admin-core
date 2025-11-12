package token

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/tokenservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTokenLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取令牌列表
func NewListTokenLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListTokenLogic {
	return &ListTokenLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListTokenLogic) ListToken(req *types.TokenListRequest) (resp *types.TokenListResponse, err error) {
	// 构建 RPC 请求
	rpcReq := &tokenservice.TokenListRequest{
		Page: &tokenservice.BasePageRequest{
			PageNumber: req.Page.CurrentPage,
			PageSize:   req.Page.PageSize,
		},
		UserId:     req.UserId,
		TokenType:  req.TokenType,
		IpAddress:  req.IpAddress,
		DeviceInfo: req.DeviceInfo,
		ProviderId: req.ProviderId,
	}

	// 调用 RPC 服务获取 Token 列表
	rpcResp, err := l.svcCtx.TokenRpc.ListToken(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	tokenList := make([]types.TokenInfo, 0, len(rpcResp.List))
	for _, token := range rpcResp.List {
		tokenList = append(tokenList, types.TokenInfo{
			Id:           token.Id,
			CreatedAt:    token.CreatedAt,
			UpdatedAt:    token.UpdatedAt,
			State:        token.State,
			TokenValue:   token.TokenValue,
			TokenType:    token.TokenType,
			UserId:       token.UserId,
			ExpiresAt:    token.ExpiresAt,
			DeviceInfo:   token.DeviceInfo,
			IpAddress:    token.IpAddress,
			LastUsedAt:   token.LastUsedAt,
			UserAgent:    token.UserAgent,
			Metadata:     token.Metadata,
			ProviderId:   token.ProviderId,
			Username:     token.Username,
			ProviderName: token.ProviderName,
			FullName:     token.FullName,
		})
	}

	resp = &types.TokenListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.TokenInfoList{
			List: tokenList,
			BaseListInfo: types.BaseListInfo{
				Total: rpcResp.Page.Total,
			},
		},
	}
	return resp, nil
}
