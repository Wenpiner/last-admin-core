package positionservicelogic

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

type CreateOrUpdatePositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdatePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdatePositionLogic {
	return &CreateOrUpdatePositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建或更新岗位
func (l *CreateOrUpdatePositionLogic) CreateOrUpdatePosition(in *core.PositionInfo) (*core.PositionInfo, error) {
	// 开启事务，并先进行检查是否存在，如果存在则进行更新否则进行创建
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	var position *ent.Position

	if in.Id != nil && *in.Id != 0 {
		// 更新操作
		if err := l.validateUpdate(in); err != nil {
			return nil, err
		}

		updateQuery := tx.Position.UpdateOneID(*in.Id)

		// 设置可更新的字段
		if in.PositionName != nil {
			updateQuery.SetPositionName(*in.PositionName)
		}
		if in.PositionCode != nil {
			updateQuery.SetPositionCode(*in.PositionCode)
		}
		if in.SortOrder != nil {
			updateQuery.SetSort(*in.SortOrder)
		}
		if in.State != nil {
			updateQuery.SetState(*in.State)
		}
		if in.Description != nil {
			updateQuery.SetNillableDescription(in.Description)
		}

		position, err = updateQuery.Save(l.ctx)
	} else {
		// 创建操作
		if err := l.validateCreate(in); err != nil {
			return nil, err
		}

		createQuery := tx.Position.Create().
			SetPositionName(*in.PositionName).
			SetPositionCode(*in.PositionCode).
			SetSort(l.getSortOrderValue(in.SortOrder)).
			SetState(l.getStateValue(in.State)).
			SetNillableDescription(in.Description)

		position, err = createQuery.Save(l.ctx)
	}

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertPositionToPositionInfo(position), nil
}

// 验证新增参数可用性
func (l *CreateOrUpdatePositionLogic) validateCreate(in *core.PositionInfo) error {
	if in.PositionName == nil || in.PositionCode == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 验证更新参数可用性
func (l *CreateOrUpdatePositionLogic) validateUpdate(in *core.PositionInfo) error {
	if in.Id == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 获取排序值，默认为 0
func (l *CreateOrUpdatePositionLogic) getSortOrderValue(sortOrder *int32) int32 {
	if sortOrder != nil {
		return *sortOrder
	}
	return 0
}

// 获取状态值，默认为 true (启用)
func (l *CreateOrUpdatePositionLogic) getStateValue(state *bool) bool {
	if state != nil {
		return *state
	}
	return true
}

// 将 Position 实体转换为 PositionInfo
func (l *CreateOrUpdatePositionLogic) convertPositionToPositionInfo(pos *ent.Position) *core.PositionInfo {
	return &core.PositionInfo{
		Id:           &pos.ID,
		CreatedAt:    pointer.ToInt64Ptr(pos.CreatedAt.UnixMilli()),
		UpdatedAt:    pointer.ToInt64Ptr(pos.UpdatedAt.UnixMilli()),
		PositionName: &pos.PositionName,
		PositionCode: &pos.PositionCode,
		SortOrder:    &pos.Sort,
		State:        &pos.State,
		Description:  pos.Description,
	}
}
