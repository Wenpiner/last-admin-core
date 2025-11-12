package token

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/tokenservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CleanExpiredTokensLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 清理过期令牌
func NewCleanExpiredTokensLogic(r *http.Request, svcCtx *svc.ServiceContext) *CleanExpiredTokensLogic {
	return &CleanExpiredTokensLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CleanExpiredTokensLogic) CleanExpiredTokens(req *types.CleanExpiredTokensRequest) (resp *types.CleanExpiredTokensResponse, err error) {
	// 调用 RPC 服务清理过期 Token
	rpcReq := &tokenservice.CleanExpiredTokensRequest{
		TokenType:  req.TokenType,
		BeforeTime: req.BeforeTime,
	}
	rpcResp, err := l.svcCtx.TokenRpc.CleanExpiredTokens(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	resp = &types.CleanExpiredTokensResponse{
		Message:      rpcResp.Message,
		CleanedCount: rpcResp.CleanedCount,
	}
	return resp, nil
}
