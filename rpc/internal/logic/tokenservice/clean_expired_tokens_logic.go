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

type CleanExpiredTokensLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCleanExpiredTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanExpiredTokensLogic {
	return &CleanExpiredTokensLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 清理过期Token
func (l *CleanExpiredTokensLogic) CleanExpiredTokens(in *core.CleanExpiredTokensRequest) (*core.CleanExpiredTokensResponse, error) {
	// 构建查询条件
	query := l.svcCtx.DBEnt.Token.Delete()

	// 如果指定了token类型，添加类型过滤
	if in.TokenType != nil && *in.TokenType != "" {
		query = query.Where(token.TokenTypeEQ(*in.TokenType))
	}

	// 设置时间条件
	var beforeTime time.Time
	if in.BeforeTime != nil && *in.BeforeTime > 0 {
		// 使用指定的时间
		beforeTime = time.Unix(*in.BeforeTime, 0)
	} else {
		// 默认清理当前时间之前的过期token
		beforeTime = time.Now()
	}

	// 添加过期时间条件
	query = query.Where(token.ExpiresAtLT(beforeTime))

	// 执行删除
	affected, err := query.Exec(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return &core.CleanExpiredTokensResponse{
		CleanedCount: int64(affected),
		Message:      "Expired tokens cleaned successfully",
	}, nil
}
