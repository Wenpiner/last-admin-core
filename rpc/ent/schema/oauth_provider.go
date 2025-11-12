package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/wenpiner/last-admin-common/ent/mixins"
)

// OauthProvider holds the schema definition for the OauthProvider entity.
type OauthProvider struct {
	ent.Schema
}

// Mixin of the OauthProvider.
func (OauthProvider) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID32Mixin{},
		mixins.TimestampMixin{},
		mixins.StateMixin{},
	}
}

// Fields of the OauthProvider.
func (OauthProvider) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider_name").
			MaxLen(100).
			NotEmpty().
			Comment("提供商名称 / Provider name"),
		field.String("provider_code").
			MaxLen(50).
			NotEmpty().
			Comment("提供商编码 / Provider code"),
		field.String("client_id").
			MaxLen(255).
			NotEmpty().
			Comment("客户端ID / Client ID"),
		field.String("client_secret").
			MaxLen(255).
			NotEmpty().
			Comment("客户端密钥 / Client secret"),
		field.String("redirect_uri").
			MaxLen(500).
			NotEmpty().
			Comment("重定向URI / Redirect URI"),
		field.String("scopes").
			MaxLen(255).
			Optional().
			Nillable().
			Comment("授权范围 / Scopes"),
		field.String("authorization_url").
			MaxLen(500).
			NotEmpty().
			Comment("授权URL / Authorization URL"),
		field.String("token_url").
			MaxLen(500).
			NotEmpty().
			Comment("令牌URL / Token URL"),
		field.String("userinfo_url").
			MaxLen(500).
			Optional().
			Nillable().
			Comment("用户信息URL / User info URL"),
		field.String("logout_url").
			MaxLen(500).
			Optional().
			Nillable().
			Comment("登出URL / Logout URL"),
		field.Uint8("auth_style").
			Default(0).
			Comment("认证方式 / Auth style"),
	}
}

// Edges of the OauthProvider.
func (OauthProvider) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tokens", Token.Type).Ref("provider"),
	}
}

// Indexes of the OauthProvider.
func (OauthProvider) Indexes() []ent.Index {
	return []ent.Index{
		// 唯一索引：提供商编码
		index.Fields("provider_code").
			Unique().
			StorageKey("sys_oauth_providers_code_unique"),
	}
}

// Annotations of the OauthProvider.
func (OauthProvider) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_oauth_providers"},
		entsql.WithComments(true),
		schema.Comment("OAuth提供商表 / OAuth provider table"),
	}
}
