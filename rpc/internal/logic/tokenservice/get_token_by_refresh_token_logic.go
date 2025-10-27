package tokenservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenByRefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenByRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByRefreshTokenLogic {
	return &GetTokenByRefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据刷新Token ID获取相关的访问Token
func (l *GetTokenByRefreshTokenLogic) GetTokenByRefreshToken(in *core.GetTokenByRefreshTokenRequest) (*core.TokenInfo, error) {
	// todo: add your logic here and delete this line

	return &core.TokenInfo{}, nil
}
