package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/wenpiner/last-admin-common/ent/mixins"
)

type Configuration struct {
	ent.Schema
}

func (Configuration) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("配置名称 / Configuration name"),
		field.String("group").NotEmpty().Comment("配置分组 / Configuration group"),
		field.String("key").NotEmpty().Comment("配置键 / Configuration key"),
		field.Text("value").NotEmpty().Comment("配置值 / Configuration value"),
		field.String("description").Optional().Comment("配置描述 / Configuration description"),
	}
}
func (Configuration) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StateMixin{},
	}
}

func (Configuration) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key").Unique(),
	}
}

func (Configuration) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "sys_configuration"},
		schema.Comment("配置表 / Configuration table"),
	}
}
