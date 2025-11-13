package configurationservicelogic

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/wenpiner/last-admin-core/rpc/ent/configuration"
	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListConfigurationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListConfigurationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConfigurationLogic {
	return &ListConfigurationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取配置列表
func (l *ListConfigurationLogic) ListConfiguration(in *core.ConfigurationListRequest) (*core.ConfigurationListResponse, error) {
	// 构建查询条件
	var predicates []predicate.Configuration

	// 根据 key 模糊搜索
	if in.Key != nil && *in.Key != "" {
		predicates = append(predicates, configuration.KeyContains(*in.Key))
	}

	// 根据 name 模糊搜索
	if in.Name != nil && *in.Name != "" {
		predicates = append(predicates, configuration.NameContains(*in.Name))
	}

	// 获取当前角色允许的读权限分组列表
	permChecker := NewConfigurationPermissionChecker(l.svcCtx.Casbin, l.Logger)
	allowedGroups, err := permChecker.GetAllowedGroups(l.ctx, OperationRead)
	if err != nil {
		return nil, err
	}

	if len(allowedGroups) == 0 {
		return &core.ConfigurationListResponse{
			Page: &core.BasePageResp{
				Total:      0,
				PageNumber: in.Page.PageNumber,
				PageSize:   in.Page.PageSize,
			},
		}, nil
	}

	// 如果指定了 group 参数，需要检查权限
	if in.Group != nil && *in.Group != "" {
		// 检查指定的 group 是否在允许列表中
		hasPermission := false
		for _, allowedGroup := range allowedGroups {
			if allowedGroup == *in.Group {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			return nil, errorx.NewInvalidArgumentError("common.configuration.forbidden")
		}
		// 精确匹配指定的 group
		predicates = append(predicates, configuration.GroupEQ(*in.Group))
	} else {
		// 没有指定 group，则只查询允许的分组
		predicates = append(predicates, configuration.GroupIn(allowedGroups...))
	}

	// 根据 key 模糊搜索
	if in.Key != nil && *in.Key != "" {
		predicates = append(predicates, configuration.KeyContains(*in.Key))
	}

	// 根据 name 模糊搜索
	if in.Name != nil && *in.Name != "" {
		predicates = append(predicates, configuration.NameContains(*in.Name))
	}

	// 执行分页查询
	page, err := l.svcCtx.DBEnt.Configuration.Query().
		Where(predicates...).
		Order(configuration.ByID(sql.OrderDesc())).
		Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 构建响应
	resp := &core.ConfigurationListResponse{
		Page: &core.BasePageResp{
			Total:      page.PageDetails.Total,
			PageNumber: in.Page.PageNumber,
			PageSize:   in.Page.PageSize,
		},
	}

	// 转换配置列表
	for _, config := range page.List {
		resp.List = append(resp.List, ConvertConfigurationToConfigurationInfo(config))
	}

	return resp, nil
}
