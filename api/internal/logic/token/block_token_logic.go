package token

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/tokenservice"
	"k8s.io/utils/pointer"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlockTokenLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 拉黑某一个令牌
func NewBlockTokenLogic(r *http.Request, svcCtx *svc.ServiceContext) *BlockTokenLogic {
	return &BlockTokenLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *BlockTokenLogic) BlockToken(req *types.BlockTokenRequest) (resp *types.BaseResponse, err error) {
	// 调用Update里面的State
	rpcReq := &tokenservice.TokenInfo{
		Id:         &req.ID,	
		State:      pointer.Bool(false),
	}
	_, err = l.svcCtx.TokenRpc.UpdateToken(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "block.token.success",
	}
	return resp, nil
}
