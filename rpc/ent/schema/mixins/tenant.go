package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type TenantMixin struct {
	mixin.Schema
}

func (TenantMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("tenant_id").
			Comment("租户ID / Tenant ID"),
	}
}
