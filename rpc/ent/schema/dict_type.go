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

// DictType holds the schema definition for the DictType entity.
type DictType struct {
	ent.Schema
}

// Mixin of the DictType.
func (DictType) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID32Mixin{},
		mixins.TimestampMixin{},
		mixins2.SoftDeleteMixin{},
		mixins.StateMixin{},
	}
}

// Fields of the DictType.
func (DictType) Fields() []ent.Field {
	return []ent.Field{
		field.String("dict_type_code").
			MaxLen(100).
			NotEmpty().
			Comment("字典类型编码 / Dictionary type code"),
		field.String("dict_type_name").
			MaxLen(100).
			NotEmpty().
			Comment("字典类型名称 / Dictionary type name"),
		field.Text("description").
			NotEmpty().
			Default("").
			Comment("字典类型描述 / Dictionary type description"),
	}
}

// Edges of the DictType.
func (DictType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dict_items", DictItem.Type).Ref("dict_type"),
	}
}

// Indexes of the DictType.
func (DictType) Indexes() []ent.Index {
	return []ent.Index{
		// 复合唯一索引：字典类型编码
		index.Fields("dict_type_code").
			Unique().
			StorageKey("sys_dict_types_code_unique"),
	}
}

// Annotations of the DictType.
func (DictType) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_dict_types"},
		entsql.WithComments(true),
		schema.Comment("字典类型表 / Dictionary type table"),
	}
}
