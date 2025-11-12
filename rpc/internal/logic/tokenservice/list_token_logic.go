package tokenservicelogic

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/ent/token"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
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
	// 构建查询条件
	var predicates []predicate.Token

	// 如果指定了用户ID，添加用户ID过滤
	if in.UserId != nil && *in.UserId != "" {
		userID, err := uuid.Parse(*in.UserId)
		if err == nil {
			predicates = append(predicates, token.UserIDEQ(userID))
		}
	}

	// 如果指定了token类型，添加类型过滤
	if in.TokenType != nil && *in.TokenType != "" {
		predicates = append(predicates, token.TokenTypeEQ(*in.TokenType))
	}

	// 如果指定了是否已撤销，添加撤销状态过滤
	if in.ProviderId != nil {
		predicates = append(predicates, token.ProviderIDEQ(*in.ProviderId))
	}

	// 如果指定了IP地址，添加IP地址过滤
	if in.IpAddress != nil && *in.IpAddress != "" {
		predicates = append(predicates, token.IPAddressEQ(*in.IpAddress))
	}

	// 如果指定了设备信息，添加设备信息过滤
	if in.DeviceInfo != nil && *in.DeviceInfo != "" {
		predicates = append(predicates, token.DeviceInfoEQ(*in.DeviceInfo))
	}

	// 按创建时间倒序排序，获取分页数据
	tokenEntities, err := l.svcCtx.DBEnt.Token.Query().
		WithUser().
		WithProvider().
		Where(predicates...).
		Order(token.ByCreatedAt(sql.OrderDesc())).
		Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 转换为TokenInfo列表
	var tokenInfos []*core.TokenInfo
	for _, tokenEntity := range tokenEntities.List {
		tokenInfos = append(tokenInfos, ConvertTokenToTokenInfo(tokenEntity))
	}

	return &core.TokenListResponse{
		Page: &core.BasePageResp{
			Total:      tokenEntities.PageDetails.Total,
			PageNumber: tokenEntities.PageDetails.Page,
			PageSize:   tokenEntities.PageDetails.Pages,
		},
		List: tokenInfos,
	}, nil
}
