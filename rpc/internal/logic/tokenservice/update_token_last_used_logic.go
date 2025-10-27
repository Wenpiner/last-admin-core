package tokenservicelogic

import (
	"context"
	"time"

	"github.com/wenpiner/last-admin-core/rpc/ent/token"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTokenLastUsedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTokenLastUsedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTokenLastUsedLogic {
	return &UpdateTokenLastUsedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新Token最后使用时间
func (l *UpdateTokenLastUsedLogic) UpdateTokenLastUsed(in *core.UpdateTokenLastUsedRequest) (*core.BaseResponse, error) {
	// 构建更新查询
	updateQuery := l.svcCtx.DBEnt.Token.Update().
		Where(token.TokenValueEQ(in.TokenValue)).
		SetLastUsedAt(time.Now())

	// 如果提供了IP地址，也更新IP地址
	if in.IpAddress != nil && *in.IpAddress != "" {
		updateQuery = updateQuery.SetIPAddress(*in.IpAddress)
	}

	// 执行更新
	affected, err := updateQuery.Save(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	if affected == 0 {
		return &core.BaseResponse{
			Message: "Token not found",
		}, nil
	}

	return &core.BaseResponse{
		Message: "Token last used time updated successfully",
	}, nil
}
