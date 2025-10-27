package oauthproviderservicelogic

import (
	"context"

	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateOauthProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdateOauthProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdateOauthProviderLogic {
	return &CreateOrUpdateOauthProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建或更新提供商
func (l *CreateOrUpdateOauthProviderLogic) CreateOrUpdateOauthProvider(in *core.OauthProviderInfo) (*core.OauthProviderInfo, error) {
	// 验证必填字段
	if err := l.validateOauthProviderInfo(in); err != nil {
		return nil, err
	}

	// 开启事务
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	var result *ent.OauthProvider

	if in.Id != nil && *in.Id > 0 {
		// 更新操作
		updateQuery := tx.OauthProvider.UpdateOneID(*in.Id)

		// 设置可更新的字段
		if in.ProviderName != nil {
			updateQuery.SetProviderName(*in.ProviderName)
		}
		if in.ProviderCode != nil {
			updateQuery.SetProviderCode(*in.ProviderCode)
		}
		if in.ClientId != nil {
			updateQuery.SetClientID(*in.ClientId)
		}
		if in.ClientSecret != nil {
			updateQuery.SetClientSecret(*in.ClientSecret)
		}
		if in.RedirectUri != nil {
			updateQuery.SetRedirectURI(*in.RedirectUri)
		}
		if in.Scopes != nil {
			updateQuery.SetNillableScopes(in.Scopes)
		}
		if in.AuthorizationUrl != nil {
			updateQuery.SetAuthorizationURL(*in.AuthorizationUrl)
		}
		if in.TokenUrl != nil {
			updateQuery.SetTokenURL(*in.TokenUrl)
		}
		if in.UserinfoUrl != nil {
			updateQuery.SetNillableUserinfoURL(in.UserinfoUrl)
		}
		if in.LogoutUrl != nil {
			updateQuery.SetNillableLogoutURL(in.LogoutUrl)
		}
		if in.AuthStyle != nil {
			updateQuery.SetAuthStyle(uint8(*in.AuthStyle))
		}
		if in.State != nil {
			updateQuery.SetState(*in.State)
		}

		result, err = updateQuery.Save(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
	} else {
		// 创建操作
		createQuery := tx.OauthProvider.Create().
			SetProviderName(*in.ProviderName).
			SetProviderCode(*in.ProviderCode).
			SetClientID(*in.ClientId).
			SetClientSecret(*in.ClientSecret).
			SetRedirectURI(*in.RedirectUri).
			SetAuthorizationURL(*in.AuthorizationUrl).
			SetTokenURL(*in.TokenUrl).
			SetAuthStyle(l.getAuthStyleValue(in.AuthStyle)).
			SetState(l.getStateValue(in.State))

		// 设置可选字段
		if in.Scopes != nil {
			createQuery.SetNillableScopes(in.Scopes)
		}
		if in.UserinfoUrl != nil {
			createQuery.SetNillableUserinfoURL(in.UserinfoUrl)
		}
		if in.LogoutUrl != nil {
			createQuery.SetNillableLogoutURL(in.LogoutUrl)
		}

		result, err = createQuery.Save(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertOauthProviderToOauthProviderInfo(result), nil
}

// 验证OAuth提供商信息的必填字段
func (l *CreateOrUpdateOauthProviderLogic) validateOauthProviderInfo(in *core.OauthProviderInfo) error {
	// 创建时的必填字段验证
	if in.Id == nil || *in.Id == 0 {
		if in.ProviderName == nil || *in.ProviderName == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.ProviderCode == nil || *in.ProviderCode == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.ClientId == nil || *in.ClientId == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.ClientSecret == nil || *in.ClientSecret == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.RedirectUri == nil || *in.RedirectUri == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.AuthorizationUrl == nil || *in.AuthorizationUrl == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.TokenUrl == nil || *in.TokenUrl == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
	}
	return nil
}

// 获取认证方式值，默认为 0
func (l *CreateOrUpdateOauthProviderLogic) getAuthStyleValue(value *uint32) uint8 {
	if value != nil {
		return uint8(*value)
	}
	return 0
}

// 获取状态值，默认为 true
func (l *CreateOrUpdateOauthProviderLogic) getStateValue(value *bool) bool {
	if value != nil {
		return *value
	}
	return true
}

// 将 OauthProvider 实体转换为 OauthProviderInfo
func (l *CreateOrUpdateOauthProviderLogic) convertOauthProviderToOauthProviderInfo(provider *ent.OauthProvider) *core.OauthProviderInfo {
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
