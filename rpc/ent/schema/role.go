package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/wenpiner/last-admin-common/ent/mixins"
	mixins2 "github.com/wenpiner/last-admin-core/rpc/ent/schema/mixins"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("role_name").
			MaxLen(50).
			NotEmpty().
			Comment("角色名称 / Role name"),
		field.String("role_code").
			MaxLen(50).
			NotEmpty().
			Comment("角色编码 / Role code"),
		field.Text("description").
			Optional().
			Nillable().
			Comment("角色描述 / Role description"),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("menus", Menu.Type),
		edge.To("users", User.Type),
	}
}

// Mixin of the Role.
func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimestampMixin{},
		mixins2.SoftDeleteMixin{},
		mixins.ID32Mixin{},
		mixins.StateMixin{},
	}
}

// Indexes of the Role.
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		// 唯一索引：角色编码（排除已删除的记录）
		index.Fields("role_code").
			Unique(),
	}
}

// Annotations of the Role.
func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "sys_roles",
		},
	}
}
