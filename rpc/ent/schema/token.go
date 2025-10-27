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

// Token holds the schema definition for the Token entity.
type Token struct {
	ent.Schema
}

// Mixin of the Token.
func (Token) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimestampMixin{},
		mixins.ID32Mixin{},
		mixins.StateMixin{},
	}
}

// Fields of the Token.
func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.String("token_value").
			MaxLen(500).
			NotEmpty().
			Comment("令牌值 / Token value"),
		field.String("token_type").
			MaxLen(50).
			NotEmpty().
			Comment("令牌类型 / Token type (access_token, refresh_token, reset_password, email_verify, api_token, sso_token)"),
		field.UUID("user_id", uuid.UUID{}).
			Optional().
			Nillable().
			Comment("用户ID / User ID"),
		field.Time("expires_at").
			Comment("过期时间 / Expiration time"),
		field.Bool("is_revoked").
			Default(false).
			Comment("是否已撤销 / Whether revoked"),
		field.String("device_info").
			MaxLen(500).
			Optional().
			Nillable().
			Comment("设备信息 / Device information"),
		field.String("ip_address").
			MaxLen(45).
			Optional().
			Nillable().
			Comment("创建时IP地址 / IP address when created"),
		field.Time("last_used_at").
			Optional().
			Nillable().
			Comment("最后使用时间 / Last used time"),
		field.String("user_agent").
			MaxLen(1000).
			Optional().
			Nillable().
			Comment("用户代理 / User agent"),
		field.Text("metadata").
			Optional().
			Nillable().
			Comment("元数据 / Metadata (JSON format)"),
		field.String("refresh_token_id").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("关联的刷新令牌ID / Associated refresh token ID"),
	}
}

// Edges of the Token.
func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Field("user_id").
			Unique(),
	}
}

// Indexes of the Token.
func (Token) Indexes() []ent.Index {
	return []ent.Index{
		// 唯一索引：令牌值
		index.Fields("token_value").
			Unique().
			StorageKey("sys_tokens_token_value_unique"),
		// 普通索引：用户ID
		index.Fields("user_id"),
		// 普通索引：令牌类型
		index.Fields("token_type"),
		// 普通索引：过期时间
		index.Fields("expires_at"),
		// 复合索引：用户ID+令牌类型+是否撤销
		index.Fields("user_id", "token_type", "is_revoked"),
		// 普通索引：刷新令牌ID
		index.Fields("refresh_token_id"),
	}
}

// Annotations of the Token.
func (Token) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_tokens"},
		entsql.WithComments(true),
		schema.Comment("令牌表 / Token table"),
	}
}
