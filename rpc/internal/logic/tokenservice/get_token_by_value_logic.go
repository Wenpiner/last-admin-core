package tokenservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/token"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenByValueLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenByValueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByValueLogic {
	return &GetTokenByValueLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据Token值获取Token信息
func (l *GetTokenByValueLogic) GetTokenByValue(in *core.StringRequest) (*core.TokenInfo, error) {
	// 查询Token
	tokenEntity, err := l.svcCtx.DBEnt.Token.Query().
		Where(token.TokenValueEQ(in.Value)).
		First(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return ConvertTokenToTokenInfo(tokenEntity), nil
}
