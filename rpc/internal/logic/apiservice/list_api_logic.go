package apiservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/api"
	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApiLogic {
	return &ListApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取API列表
func (l *ListApiLogic) ListApi(in *core.ApiListRequest) (*core.ApiListResponse, error) {
	// 构建查询条件
	var predicates []predicate.API

	// 服务名称过滤
	if in.ServiceName != nil && *in.ServiceName != "" {
		predicates = append(predicates, api.ServiceNameContains(*in.ServiceName))
	}

	// API分组过滤
	if in.ApiGroup != nil && *in.ApiGroup != "" {
		predicates = append(predicates, api.APIGroupContains(*in.ApiGroup))
	}

	// 请求方法过滤
	if in.Method != nil && *in.Method != "" {
		predicates = append(predicates, api.MethodEQ(*in.Method))
	}

	// 描述过滤
	if in.Description != nil && *in.Description != "" {
		predicates = append(predicates, api.DescriptionContains(*in.Description))
	}

	if in.Path != nil && *in.Path != "" {
		predicates = append(predicates, api.PathContains(*in.Path))
	}

	// 执行分页查询
	query := l.svcCtx.DBEnt.API.Query().Where(predicates...)

	// 获取总数
	total, err := query.Count(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 分页查询
	apis, err := query.
		Offset(int((in.Page.PageNumber - 1) * in.Page.PageSize)).
		Limit(int(in.Page.PageSize)).
		All(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 转换结果
	var apiList []*core.ApiInfo
	for _, apiEntity := range apis {
		apiList = append(apiList, l.convertAPIToApiInfo(apiEntity))
	}

	return &core.ApiListResponse{
		Page: &core.BasePageResp{
			Total:      uint64(total),
			PageNumber: in.Page.PageNumber,
			PageSize:   in.Page.PageSize,
		},
		List: apiList,
	}, nil
}

// 将 API 实体转换为 ApiInfo
func (l *ListApiLogic) convertAPIToApiInfo(apiEntity *ent.API) *core.ApiInfo {
	return &core.ApiInfo{
		Id:          pointer.ToUint32Ptr(uint32(apiEntity.ID)),
		CreatedAt:   pointer.ToInt64Ptr(apiEntity.CreatedAt.UnixMilli()),
		UpdatedAt:   pointer.ToInt64Ptr(apiEntity.UpdatedAt.UnixMilli()),
		Name:        pointer.ToStringPtrIfNotEmpty(pointer.GetString(apiEntity.Name)),
		Method:      &apiEntity.Method,
		Path:        &apiEntity.Path,
		Description: pointer.ToStringPtrIfNotEmpty(pointer.GetString(apiEntity.Description)),
		IsRequired:  &apiEntity.IsRequired,
		ServiceName: &apiEntity.ServiceName,
		ApiGroup:    &apiEntity.APIGroup,
	}
}
