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

type UnblockTokenLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 解封某一个令牌
func NewUnblockTokenLogic(r *http.Request, svcCtx *svc.ServiceContext) *UnblockTokenLogic {
	return &UnblockTokenLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *UnblockTokenLogic) UnblockToken(req *types.BlockTokenRequest) (resp *types.BaseResponse, err error) {
	_, err = l.svcCtx.TokenRpc.UpdateToken(l.ctx, &tokenservice.TokenInfo{
		Id:         &req.ID,
		State:      pointer.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "unblock.token.success",
	}

	return
}
