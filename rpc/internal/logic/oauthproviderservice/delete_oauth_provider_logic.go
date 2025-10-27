package oauthproviderservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/oauthprovider"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOauthProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteOauthProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOauthProviderLogic {
	return &DeleteOauthProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除提供商
func (l *DeleteOauthProviderLogic) DeleteOauthProvider(in *core.ID32Request) (*core.BaseResponse, error) {
	// 执行删除操作
	_, err := l.svcCtx.DBEnt.OauthProvider.Delete().Where(
		oauthprovider.IDEQ(in.Id),
	).Exec(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return &core.BaseResponse{
		Message: "删除成功",
	}, nil
}
