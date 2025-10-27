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
)

// OperationLog holds the schema definition for the OperationLog entity.
type OperationLog struct {
	ent.Schema
}

// Mixin of the OperationLog.
func (OperationLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID32Mixin{},
		mixins.TimestampMixin{},
	}
}

// Fields of the OperationLog.
func (OperationLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}).
			Optional().
			Nillable().
			Comment("操作用户ID / Operation user ID"),
		field.String("username").
			MaxLen(50).
			Optional().
			Nillable().
			Comment("操作用户名 / Operation username"),
		field.String("operation_type").
			MaxLen(20).
			NotEmpty().
			Comment("操作类型 / Operation type (CREATE, UPDATE, DELETE, LOGIN, LOGOUT, etc.)"),
		field.String("module").
			MaxLen(50).
			NotEmpty().
			Comment("操作模块 / Operation module"),
		field.String("business_type").
			MaxLen(50).
			Optional().
			Nillable().
			Comment("业务类型 / Business type"),
		field.String("method").
			MaxLen(10).
			Optional().
			Nillable().
			Comment("请求方法 / Request method"),
		field.String("request_url").
			MaxLen(500).
			Optional().
			Nillable().
			Comment("请求URL / Request URL"),
		field.Text("request_params").
			Optional().
			Nillable().
			Comment("请求参数 / Request parameters"),
		field.Text("response_data").
			Optional().
			Nillable().
			Comment("响应数据 / Response data"),
		field.String("ip_address").
			MaxLen(45).
			Optional().
			Nillable().
			Comment("IP地址 / IP address"),
		field.String("user_agent").
			MaxLen(500).
			Optional().
			Nillable().
			Comment("用户代理 / User agent"),
		field.String("location").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("操作地点 / Operation location"),
		field.Bool("is_success").
			Default(true).
			Comment("是否成功 / Whether successful"),
		field.String("error_message").
			MaxLen(1000).
			Optional().
			Nillable().
			Comment("错误信息 / Error message"),
		field.Int("execution_time").
			Default(0).
			Comment("执行时间(毫秒) / Execution time (milliseconds)"),
		field.Text("description").
			Optional().
			Nillable().
			Comment("操作描述 / Operation description"),
	}
}

// Edges of the OperationLog.
func (OperationLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Field("user_id").
			Unique(),
	}
}

// Indexes of the OperationLog.
func (OperationLog) Indexes() []ent.Index {
	return []ent.Index{
		// 普通索引：用户ID
		index.Fields("user_id"),
		// 普通索引：操作类型
		index.Fields("operation_type"),
		// 普通索引：模块
		index.Fields("module"),
		// 普通索引：创建时间
		index.Fields("created_at"),
	}
}

// Annotations of the OperationLog.
func (OperationLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_operation_logs"},
		entsql.WithComments(true),
		schema.Comment("操作记录表 / Operation log table"),
	}
}
