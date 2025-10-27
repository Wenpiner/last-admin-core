package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type DictType struct {
	ent.Schema
}

func (DictType) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Comment("字典类型ID / Dictionary type ID").
			Positive().
			Immutable(),
		field.String("dict_type_code").
			MaxLen(100).
			NotEmpty().
			Comment("字典类型编码 / Dictionary type code"),
		field.String("dict_type_name").
			MaxLen(100).
			NotEmpty().
			Comment("字典类型名称 / Dictionary type name"),
		field.Int8("status").
			Default(1).
			Comment("状态(1:启用,0:禁用) / Status (1:enabled, 0:disabled)"),
		field.Text("description").
			Optional().
			Comment("字典类型描述 / Dictionary type description"),
	}
}
