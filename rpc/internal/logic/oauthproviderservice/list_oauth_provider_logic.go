package oauthproviderservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/oauthprovider"
	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOauthProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOauthProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOauthProviderLogic {
	return &ListOauthProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取提供商列表
func (l *ListOauthProviderLogic) ListOauthProvider(in *core.OauthProviderListRequest) (*core.OauthProviderListResponse, error) {
	// 构建查询条件
	var predicates []predicate.OauthProvider

	// 提供商名称过滤
	if in.ProviderName != nil && *in.ProviderName != "" {
		predicates = append(predicates, oauthprovider.ProviderNameContains(*in.ProviderName))
	}

	// 提供商编码过滤
	if in.ProviderCode != nil && *in.ProviderCode != "" {
		predicates = append(predicates, oauthprovider.ProviderCodeContains(*in.ProviderCode))
	}

	// 状态过滤
	if in.State != nil {
		predicates = append(predicates, oauthprovider.StateEQ(*in.State))
	}

	// 执行分页查询
	query := l.svcCtx.DBEnt.OauthProvider.Query().Where(predicates...)

	// 分页查询
	providers, err := query.Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 转换结果
	var providerList []*core.OauthProviderInfo
	for _, provider := range providers.List {
		providerList = append(providerList, l.convertOauthProviderToOauthProviderInfo(provider))
	}

	return &core.OauthProviderListResponse{
		Page: &core.BasePageResp{
			Total:      providers.PageDetails.Total,
			PageNumber: in.Page.PageNumber,
			PageSize:   in.Page.PageSize,
		},
		List: providerList,
	}, nil
}

// 将 OauthProvider 实体转换为 OauthProviderInfo
func (l *ListOauthProviderLogic) convertOauthProviderToOauthProviderInfo(provider *ent.OauthProvider) *core.OauthProviderInfo {
	return &core.OauthProviderInfo{
		Id:               pointer.ToUint32Ptr(uint32(provider.ID)),
		CreatedAt:        pointer.ToInt64Ptr(provider.CreatedAt.UnixMilli()),
		UpdatedAt:        pointer.ToInt64Ptr(provider.UpdatedAt.UnixMilli()),
		ProviderName:     &provider.ProviderName,
		ProviderCode:     &provider.ProviderCode,
		ClientId:         &provider.ClientID,
		ClientSecret:     &provider.ClientSecret,
		RedirectUri:      &provider.RedirectURI,
		Scopes:           pointer.ToStringPtrIfNotEmpty(pointer.GetString(provider.Scopes)),
		AuthorizationUrl: &provider.AuthorizationURL,
		TokenUrl:         &provider.TokenURL,
		UserinfoUrl:      pointer.ToStringPtrIfNotEmpty(pointer.GetString(provider.UserinfoURL)),
		LogoutUrl:        pointer.ToStringPtrIfNotEmpty(pointer.GetString(provider.LogoutURL)),
		AuthStyle:        pointer.ToUint32Ptr(uint32(provider.AuthStyle)),
		State:            &provider.State,
	}
}
