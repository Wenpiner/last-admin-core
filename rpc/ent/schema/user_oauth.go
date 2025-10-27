package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/wenpiner/last-admin-common/ent/mixins"
)

// UserOauth holds the schema definition for the UserOauth entity.
type UserOauth struct {
	ent.Schema
}

// Mixin of the UserOauth.
func (UserOauth) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimestampMixin{},
		mixins.ID32Mixin{},
	}
}

// Indexes of the UserOauth.
func (UserOauth) Indexes() []ent.Index {
	return []ent.Index{
		// 复合唯一索引：认证类型+第三方用户ID
		index.Fields("provider_id", "oauth_id").
			Unique(),
	}
}

// Fields of the UserOauth.
func (UserOauth) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}).
			Comment("用户ID / User ID"),
		field.Uint32("provider_id").
			Comment("提供商ID / Provider ID"),
		field.String("oauth_id").
			MaxLen(100).
			NotEmpty().
			Comment("第三方用户ID / Third-party user ID"),
		field.JSON("oauth_data", map[string]interface{}{}).
			Optional().
			Comment("第三方用户数据 / Third-party user data"),
		field.Text("oauth_token").
			Optional().
			Nillable().
			Comment("访问令牌 / Access token"),
		field.Text("refresh_token").
			Optional().
			Nillable().
			Comment("刷新令牌 / Refresh token"),
		field.Time("token_expires_at").
			Optional().
			Nillable().
			Comment("令牌过期时间 / Token expiration time"),
	}
}

// Edges of the UserOauth.
func (UserOauth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Field("user_id").
			Unique().Required(),
		edge.To("provider", OauthProvider.Type).
			Field("provider_id").
			Unique().Required(),
	}
}

// Annotations of the UserOauth.
func (UserOauth) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "sys_user_oauth",
		},
		entsql.WithComments(true),
		schema.Comment("用户第三方认证表 / User third-party authentication table"),
	}
}
