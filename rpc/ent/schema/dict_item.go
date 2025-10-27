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

// DictItem holds the schema definition for the DictItem entity.
type DictItem struct {
	ent.Schema
}

// Mixin of the DictItem.
func (DictItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID32Mixin{},
		mixins.TimestampMixin{},
		mixins2.SoftDeleteMixin{},
		mixins.StateMixin{},
	}
}

// Fields of the DictItem.
func (DictItem) Fields() []ent.Field {
	return []ent.Field{
		field.String("item_label").
			MaxLen(100).
			NotEmpty().
			Comment("字典项标签 / Dictionary item label"),
		field.String("item_value").
			MaxLen(100).
			NotEmpty().
			Comment("字典项值 / Dictionary item value"),
		field.String("item_color").
			MaxLen(20).
			Optional().
			Nillable().
			Comment("字典项颜色 / Dictionary item color"),
		field.String("item_css").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("字典项CSS / Dictionary item CSS"),
		field.Int("sort_order").
			Default(0).
			Comment("排序 / Sort order"),
		field.Text("description").
			Optional().
			Nillable().
			Comment("字典项描述 / Dictionary item description"),
		field.Uint32("dict_type_id").
			Comment("字典类型ID / Dictionary type ID"),
	}
}

// Edges of the DictItem.
func (DictItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("dict_type", DictType.Type).
			Field("dict_type_id").
			Required().
			Unique(),
	}
}

// Indexes of the DictItem.
func (DictItem) Indexes() []ent.Index {
	return []ent.Index{
		// 复合唯一索引：字典类型ID+字典项值
		index.Fields("dict_type_id", "item_value").
			Unique().
			StorageKey("sys_dict_items_type_value_unique"),
		// 普通索引：字典类型ID
		index.Fields("dict_type_id"),
	}
}

// Annotations of the DictItem.
func (DictItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_dict_items"},
		entsql.WithComments(true),
		schema.Comment("字典项表 / Dictionary item table"),
	}
}
