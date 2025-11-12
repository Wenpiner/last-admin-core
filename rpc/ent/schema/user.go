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

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimestampMixin{},
		mixins.StateMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Comment("用户ID / User ID").StorageKey("id"),
		field.String("username").
			MaxLen(50).
			NotEmpty().
			Comment("用户名 / Username"),
		field.String("password_hash").
			MaxLen(255).
			NotEmpty().
			Comment("密码哈希 / Password hash"),
		field.String("email").
			MaxLen(100).
			Optional().
			Comment("邮箱 / Email"),
		field.String("full_name").
			MaxLen(100).
			Optional().
			Comment("全名 / Full name"),
		field.String("mobile").
			MaxLen(20).
			Optional().
			Comment("手机号 / Mobile phone"),
		field.String("avatar").
			MaxLen(255).
			Optional().
			Comment("头像 / Avatar"),
		field.Text("user_description").
			Optional().
			Comment("用户描述 / User description"),
		field.Time("last_login_at").
			Optional().
			Nillable().
			Comment("最后登录时间 / Last login time"),
		field.String("last_login_ip").
			MaxLen(45).
			Optional().
			Comment("最后登录IP / Last login IP"),
		field.Uint32("department_id").
			Optional().
			Comment("部门ID / Department ID"),
		field.String("home_path").
			MaxLen(255).
			Optional().
			Comment("首页路径 / Home path"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("roles", Role.Type).Ref("users"),
		edge.From("positions", Position.Type).Ref("users"),
		edge.To("department", Department.Type).
			Field("department_id").
			Unique(),
		edge.From("leader_department", Department.Type).Ref("leader"),

		edge.To("totp", UserTotp.Type).
			Unique(),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// 复合唯一索引：用户名（排除已删除的记录）
		index.Fields("username").
			Unique().
			StorageKey("sys_users_username_unique"),
		// 普通索引：部门ID
		index.Fields("department_id"),
	}
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_users"},
		entsql.WithComments(true),
		schema.Comment("用户表 / User table"),
	}
}
