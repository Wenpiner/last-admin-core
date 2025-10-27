package errorhandler

import (
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

// DBEntError 统一处理数据库错误
func DBEntError(logger logx.Logger, err error, in interface{}) error {
	if err == nil {
		return nil
	}

	// 记录错误日志
	logger.Errorf("Database error: %v, input: %+v", err, in)

	// 根据不同的 Ent 错误类型返回相应的业务错误
	switch {
	case ent.IsNotFound(err):
		// 数据不存在
		return errorx.NewInvalidArgumentError(last_i18n.TargetNotExist)
	case ent.IsConstraintError(err):
		// 约束错误（如唯一性约束违反）
		return errorx.NewInvalidArgumentError(last_i18n.DataConflict)
	case ent.IsValidationError(err):
		// 验证错误
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	case ent.IsNotSingular(err):
		// 数据一致性错误
		return errorx.NewInvalidArgumentError(last_i18n.ConsistencyCheckFailed)
	default:
		// 其他数据库错误
		return errorx.NewInternalError(last_i18n.DatabaseError)
	}
}
