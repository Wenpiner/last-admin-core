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

type Menu struct {
	ent.Schema
}

// Mixin of the Menu.
func (Menu) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID32Mixin{},
		mixins.TimestampMixin{},
		mixins2.SoftDeleteMixin{},
		mixins.StateMixin{},
		mixins.SortMixin{},
	}
}

// Fields of the Menu.
func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.String("menu_code").MaxLen(100).NotEmpty().Comment("菜单编码 / Menu code"),
		field.String("menu_name").MaxLen(100).NotEmpty().Comment("菜单名称 / Menu name"),
		field.Uint32("parent_id").Optional().Nillable().Comment("父菜单ID / Parent menu ID"),
		field.String("menu_path").MaxLen(255).Optional().Nillable().Comment("菜单路径 / Menu path"),
		field.String("component").MaxLen(100).Optional().Nillable().Comment("前端组件 / Frontend component"),
		field.String("redirect").MaxLen(255).Optional().Nillable().Comment("重定向地址 / Redirect path"),
		field.Uint16("menu_level").Default(0).Comment("菜单层级 / Menu level"),
		field.String("icon").MaxLen(100).Optional().Nillable().Comment("图标 / Icon"),
		field.String("permission").MaxLen(100).Optional().Nillable().Comment("权限标识 / Permission identifier"),
		field.String("service_name").MaxLen(100).Optional().Nillable().Comment("服务名称 / Service name"),
		field.String("menu_type").MaxLen(20).NotEmpty().Comment("菜单类型(menu,button) / Menu type (menu, button)"),
		field.Bool("is_hidden").Nillable().Optional().Comment("是否隐藏 / Whether it is hidden"),
		field.Bool("is_breadcrumb").Nillable().Optional().Comment("是否显示在面包屑中 / Whether to display in breadcrumb"),
		field.Bool("is_cache").Optional().Nillable().Comment("是否缓存 / Whether to cache"),
		field.Bool("is_tab").Optional().Nillable().Comment("是否显示在标签栏中 / Whether to display in tabbar"),
		field.Bool("is_affix").Optional().Nillable().Comment("是否固定在标签栏 / Whether to fix in tabbar"),
		field.String("frame_src").MaxLen(255).Optional().Nillable().Comment("内嵌iframe的url / Embedded iframe url"),
		field.Text("description").Optional().Nillable().Comment("菜单描述 / Menu description"),
		field.String("link").MaxLen(255).Optional().Nillable().Comment("外链地址 / Link address"),
	}
}

// Edges of the Menu.
func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Menu.Type).
			From("parent").
			Field("parent_id").
			Unique(),
		edge.From("roles", Role.Type).Ref("menus"),
	}
}

// Indexes of the Menu.
func (Menu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("menu_code").
			Unique().
			StorageKey("sys_menus_menu_code_unique"),
	}
}

// Annotations of the Menu.
func (Menu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_menus"},
		entsql.WithComments(true),
		schema.Comment("菜单资源表 / Menu resource table"),
	}
}
