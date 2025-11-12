package tokenservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/token"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTokenLogic {
	return &DeleteTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除Token
func (l *DeleteTokenLogic) DeleteToken(in *core.ID32Request) (*core.BaseResponse, error) {
	// 删除Token
	affected, err := l.svcCtx.DBEnt.Token.Delete().
		Where(token.IDEQ(in.Id)).
		Exec(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	if affected == 0 {
		return &core.BaseResponse{
			Message: "Token not found",
		}, nil
	}

	return &core.BaseResponse{
		Message: "Token deleted successfully",
	}, nil
}
