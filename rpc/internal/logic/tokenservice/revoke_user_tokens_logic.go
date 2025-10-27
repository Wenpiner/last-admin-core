package tokenservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeUserTokensLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeUserTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeUserTokensLogic {
	return &RevokeUserTokensLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 撤销用户的所有Token（除了blockUserAllToken功能）
func (l *RevokeUserTokensLogic) RevokeUserTokens(in *core.RevokeUserTokensRequest) (*core.BaseResponse, error) {
	// todo: add your logic here and delete this line

	return &core.BaseResponse{}, nil
}
