package apiservicelogic

import (
	"context"

	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdateApiLogic {
	return &CreateOrUpdateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建或更新API
func (l *CreateOrUpdateApiLogic) CreateOrUpdateApi(in *core.ApiInfo) (*core.ApiInfo, error) {
	// 验证必填字段
	if err := l.validateApiInfo(in); err != nil {
		return nil, err
	}

	// 开启事务
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	var result *ent.API

	if in.Id != nil && *in.Id > 0 {
		// 更新操作
		updateQuery := tx.API.UpdateOneID(*in.Id)

		// 设置可更新的字段
		if in.Name != nil {
			updateQuery.SetNillableName(in.Name)
		}
		if in.Method != nil {
			updateQuery.SetMethod(*in.Method)
		}
		if in.Path != nil {
			updateQuery.SetPath(*in.Path)
		}
		if in.Description != nil {
			updateQuery.SetNillableDescription(in.Description)
		}
		if in.IsRequired != nil {
			updateQuery.SetIsRequired(*in.IsRequired)
		}
		if in.ServiceName != nil {
			updateQuery.SetServiceName(*in.ServiceName)
		}
		if in.ApiGroup != nil {
			updateQuery.SetAPIGroup(*in.ApiGroup)
		}

		result, err = updateQuery.Save(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
	} else {
		// 创建操作
		createQuery := tx.API.Create().
			SetMethod(*in.Method).
			SetPath(*in.Path).
			SetIsRequired(l.getBoolValue(in.IsRequired)).
			SetServiceName(*in.ServiceName).
			SetAPIGroup(*in.ApiGroup)

		// 设置可选字段
		if in.Name != nil {
			createQuery.SetNillableName(in.Name)
		}
		if in.Description != nil {
			createQuery.SetNillableDescription(in.Description)
		}

		result, err = createQuery.Save(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertAPIToApiInfo(result), nil
}

// 验证API信息的必填字段
func (l *CreateOrUpdateApiLogic) validateApiInfo(in *core.ApiInfo) error {
	// 创建时的必填字段验证
	if in.Id == nil || *in.Id == 0 {
		if in.Method == nil || *in.Method == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.Path == nil || *in.Path == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.ServiceName == nil || *in.ServiceName == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.ApiGroup == nil || *in.ApiGroup == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
	}
	return nil
}

// 获取布尔值，默认为 false
func (l *CreateOrUpdateApiLogic) getBoolValue(value *bool) bool {
	if value != nil {
		return *value
	}
	return false
}

// 将 API 实体转换为 ApiInfo
func (l *CreateOrUpdateApiLogic) convertAPIToApiInfo(apiEntity *ent.API) *core.ApiInfo {
	return &core.ApiInfo{
		Id:          pointer.ToUint32Ptr(uint32(apiEntity.ID)),
		CreatedAt:   pointer.ToInt64Ptr(apiEntity.CreatedAt.UnixMilli()),
		UpdatedAt:   pointer.ToInt64Ptr(apiEntity.UpdatedAt.UnixMilli()),
		Name:        pointer.ToStringPtrIfNotEmpty(pointer.GetString(apiEntity.Name)),
		Method:      &apiEntity.Method,
		Path:        &apiEntity.Path,
		Description: pointer.ToStringPtrIfNotEmpty(pointer.GetString(apiEntity.Description)),
		IsRequired:  &apiEntity.IsRequired,
		ServiceName: &apiEntity.ServiceName,
		ApiGroup:    &apiEntity.APIGroup,
	}
}
