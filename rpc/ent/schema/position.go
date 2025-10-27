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

// Position holds the schema definition for the Position entity.
type Position struct {
	ent.Schema
}

// Mixin of the Position.
func (Position) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID32Mixin{},
		mixins.TimestampMixin{},
		mixins2.SoftDeleteMixin{},
		mixins.StateMixin{},
		mixins.SortMixin{},
	}
}

// Fields of the Position.
func (Position) Fields() []ent.Field {
	return []ent.Field{
		field.String("position_name").
			MaxLen(100).
			NotEmpty().
			Comment("职位名称 / Position name"),
		field.String("position_code").
			MaxLen(50).
			NotEmpty().
			Comment("职位编码 / Position code"),
		field.Text("description").
			Optional().
			Nillable().
			Comment("职位描述 / Position description"),
	}
}

// Edges of the Position.
func (Position) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}

// Indexes of the Position.
func (Position) Indexes() []ent.Index {
	return []ent.Index{
		// 复合唯一索引：职位编码
		index.Fields("position_code").
			Unique().
			StorageKey("sys_positions_code_unique"),
	}
}

// Annotations of the Position.
func (Position) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_positions"},
		entsql.WithComments(true),
		schema.Comment("职位表 / Position table"),
	}
}
