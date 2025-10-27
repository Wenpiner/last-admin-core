package tokenservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTokenLogic {
	return &ListTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取Token列表
func (l *ListTokenLogic) ListToken(in *core.TokenListRequest) (*core.TokenListResponse, error) {
	// todo: add your logic here and delete this line

	return &core.TokenListResponse{}, nil
}
