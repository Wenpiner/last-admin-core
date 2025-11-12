package tokenservicelogic

import (
	"context"

	last_redis "github.com/wenpiner/last-admin-common/last-redis"
	"github.com/wenpiner/last-admin-core/rpc/ent/token"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTokenLogic {
	return &UpdateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新Token
func (l *UpdateTokenLogic) UpdateToken(in *core.TokenInfo) (*core.TokenInfo, error) {
	// 验证ID
	if in.Id == nil || *in.Id == 0 {
		return &core.TokenInfo{}, nil
	}

	// 构建更新查询
	updateQuery := l.svcCtx.DBEnt.Token.Update().
		Where(token.IDEQ(*in.Id))

	// 更新允许的字段
	if in.State != nil {
		updateQuery = updateQuery.SetState(*in.State)
	}


	if in.DeviceInfo != nil {
		updateQuery = updateQuery.SetNillableDeviceInfo(in.DeviceInfo)
	}

	if in.IpAddress != nil {
		updateQuery = updateQuery.SetNillableIPAddress(in.IpAddress)
	}

	if in.UserAgent != nil {
		updateQuery = updateQuery.SetNillableUserAgent(in.UserAgent)
	}

	if in.Metadata != nil {
		updateQuery = updateQuery.SetNillableMetadata(in.Metadata)
	}

	// 执行更新
	affected, err := updateQuery.Save(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	
	if affected == 0 {
		return &core.TokenInfo{}, nil
	}

	if in.State != nil {
		if *in.State {
			// 从黑名单中移除
			l.svcCtx.Redis.SRem(l.ctx, string(last_redis.BlacklistToken), in.TokenValue).Result()
		} else {
			// 加入黑名单
			l.svcCtx.Redis.SAdd(l.ctx, string(last_redis.BlacklistToken), in.TokenValue).Result()
		}
	}

	// 查询更新后的Token
	tokenEntity, err := l.svcCtx.DBEnt.Token.Query().
		Where(token.IDEQ(*in.Id)).
		First(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	
	return ConvertTokenToTokenInfo(tokenEntity), nil
}
