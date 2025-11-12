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

// UserTotp holds the schema definition for the UserTotp entity.
type UserTotp struct {
	ent.Schema
}

// Mixin of the UserTotp.
func (UserTotp) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimestampMixin{},
		mixins.StateMixin{},
	}
}

// Fields of the UserTotp.
func (UserTotp) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("secret_key").
			MaxLen(255).
			NotEmpty().
			Comment("TOTP密钥 / TOTP secret key"),
		field.String("backup_codes").
			MaxLen(1000).
			Optional().
			Nillable().
			Comment("备用恢复码 / Backup recovery codes (JSON array)"),
		field.Bool("is_verified").
			Default(false).
			Comment("是否已验证 / Whether verified"),
		field.Time("last_used_at").
			Optional().
			Nillable().
			Comment("最后使用时间 / Last used time"),
		field.String("last_used_code").
			MaxLen(10).
			Optional().
			Nillable().
			Comment("最后使用的验证码 / Last used verification code"),
		field.String("device_name").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("设备名称 / Device name"),
		field.String("issuer").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("发行者名称 / Issuer name"),
	}
}

// Edges of the UserTotp.
func (UserTotp) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("totp"). // 引用 User schema 中定义的 "profile" 边
			Unique().       
			Required(),
	}
}

// Indexes of the UserTotp.
func (UserTotp) Indexes() []ent.Index {
	return []ent.Index{
		// 复合唯一索引：用户ID
		index.Fields("id").
			Unique().
			StorageKey("sys_user_totp_user_unique"),
		// 普通索引：用户ID
		index.Fields("id"),
		// 普通索引：是否已验证
		index.Fields("is_verified"),
		// 普通索引：最后使用时间
		index.Fields("last_used_at"),
	}
}

// Annotations of the UserTotp.
func (UserTotp) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_user_totp"},
		entsql.WithComments(true),
		schema.Comment("用户TOTP表 / User TOTP table"),
	}
}
