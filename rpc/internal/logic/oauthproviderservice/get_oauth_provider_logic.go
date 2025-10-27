package oauthproviderservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOauthProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOauthProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOauthProviderLogic {
	return &GetOauthProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取提供商
func (l *GetOauthProviderLogic) GetOauthProvider(in *core.ID32Request) (*core.OauthProviderInfo, error) {
	// 查询OAuth提供商
	provider, err := l.svcCtx.DBEnt.OauthProvider.Get(l.ctx, in.Id)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertOauthProviderToOauthProviderInfo(provider), nil
}

// 将 OauthProvider 实体转换为 OauthProviderInfo
func (l *GetOauthProviderLogic) convertOauthProviderToOauthProviderInfo(provider *ent.OauthProvider) *core.OauthProviderInfo {
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
