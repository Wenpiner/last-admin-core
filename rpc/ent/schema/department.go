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
	mixins2 "github.com/wenpiner/last-admin-core/rpc/ent/schema/mixins"
)

// Department holds the schema definition for the Department entity.
type Department struct {
	ent.Schema
}

// Mixin of the Department.
func (Department) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID32Mixin{},
		mixins.TimestampMixin{},
		mixins2.SoftDeleteMixin{},
		mixins.StateMixin{},
		mixins.SortMixin{},
	}
}

// Fields of the Department.
func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.String("dept_name").
			MaxLen(100).
			NotEmpty().
			Comment("部门名称 / Department name"),
		field.String("dept_code").
			MaxLen(50).
			NotEmpty().
			Comment("部门编码 / Department code"),
		field.Uint32("parent_id").
			Optional().
			Nillable().
			Comment("父部门ID / Parent department ID"),
		field.UUID("leader_user_id", uuid.UUID{}).
			Optional().
			Nillable().Unique().
			Comment("部门负责人用户ID / Leader user ID"),
		field.Text("description").
			Optional().
			Nillable().
			Comment("部门描述 / Department description"),
	}
}

// Edges of the Department.
func (Department) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Department.Type).
			From("parent").
			Field("parent_id").
			Unique(),
		edge.From("users", User.Type).Ref("department"),
		edge.To("leader", User.Type).
			Field("leader_user_id").
			Unique(),
	}
}

// Indexes of the Department.
func (Department) Indexes() []ent.Index {
	return []ent.Index{
		// 唯一索引：部门编码
		index.Fields("dept_code").
			Unique().
			StorageKey("sys_departments_code_unique"),
		// 普通索引：父部门ID
		index.Fields("parent_id"),
	}
}

// Annotations of the Department.
func (Department) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_departments"},
		entsql.WithComments(true),
		schema.Comment("部门表 / Department table"),
	}
}
