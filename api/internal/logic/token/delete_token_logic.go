package token

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/tokenservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTokenLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除令牌
func NewDeleteTokenLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteTokenLogic {
	return &DeleteTokenLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeleteTokenLogic) DeleteToken(req *types.ID32Request) (resp *types.BaseResponse, err error) {
	// 调用 RPC 服务删除 Token
	rpcReq := &tokenservice.ID32Request{
		Id: req.ID,
	}
	rpcResp, err := l.svcCtx.TokenRpc.DeleteToken(l.ctx, rpcReq)
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
