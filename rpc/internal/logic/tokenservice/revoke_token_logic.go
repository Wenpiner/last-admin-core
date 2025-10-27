package tokenservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/token"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeTokenLogic {
	return &RevokeTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 撤销Token
func (l *RevokeTokenLogic) RevokeToken(in *core.RevokeTokenRequest) (*core.BaseResponse, error) {
	// 更新Token状态为已撤销
	affected, err := l.svcCtx.DBEnt.Token.Update().
		Where(token.TokenValueEQ(in.TokenValue)).
		SetIsRevoked(true).
		Save(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	if affected == 0 {
		return &core.BaseResponse{
			Message: "Token not found or already revoked",
		}, nil
	}

	return &core.BaseResponse{
		Message: "Token revoked successfully",
	}, nil
}
