package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/wenpiner/last-admin-common/ent/mixins"
)

type API struct {
	ent.Schema
}

func (API) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("API名称 / API name"),
		field.String("method").
			MaxLen(10).
			NotEmpty().
			Comment("请求方法 / Request method"),
		field.String("path").
			MaxLen(255).
			NotEmpty().
			Comment("请求路径 / Request path"),
		field.Text("description").
			Optional().
			Nillable().
			Comment("API描述 / API description"),
		field.Bool("is_required").
			Default(false).
			Comment("是否必填 / Whether it is required"),
		field.String("service_name").
			MaxLen(100).
			NotEmpty().
			Comment("服务名称 / Service name"),
		field.String("api_group").
			MaxLen(100).
			NotEmpty().
			Comment("API分组 / API group"),
	}
}

func (API) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID32Mixin{},
		mixins.TimestampMixin{},
	}
}

func (API) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("path", "method").
			Unique(),
	}
}

func (API) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "sys_apis"},
		schema.Comment("API表 / API table"),
	}
}
