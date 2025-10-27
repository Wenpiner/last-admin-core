package dictservicelogic

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

type CreateOrUpdateDictItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdateDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdateDictItemLogic {
	return &CreateOrUpdateDictItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建或更新字典子项
func (l *CreateOrUpdateDictItemLogic) CreateOrUpdateDictItem(in *core.DictItemInfo) (*core.DictItemInfo, error) {
	// 开启事务，并先进行检查是否存在，如果存在则进行更新否则进行创建
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	var dictItem *ent.DictItem
	if in.Id == nil {
		// 新增,验证必填参数可用性
		if err := l.validateCreate(in); err != nil {
			return nil, err
		}
		sortOrder := 0
		if in.SortOrder != nil {
			sortOrder = int(*in.SortOrder)
		}
		dictItem, err = tx.DictItem.Create().
			SetItemLabel(pointer.GetString(in.Label)).
			SetItemValue(pointer.GetString(in.Value)).
			SetNillableItemColor(in.Color).
			SetNillableItemCSS(in.Css).
			SetSortOrder(sortOrder).
			SetNillableDescription(in.Description).
			SetState(pointer.GetBool(in.State)).
			SetDictTypeID(pointer.GetUint32(in.DictTypeId)).
			Save(l.ctx)
	} else {
		// 更新
		if err := l.validateUpdate(in); err != nil {
			return nil, err
		}
		var sortOrderPtr *int
		if in.SortOrder != nil {
			sortOrder := int(*in.SortOrder)
			sortOrderPtr = &sortOrder
		}
		dictItem, err = tx.DictItem.UpdateOneID(pointer.GetUint32(in.Id)).
			SetNillableItemLabel(in.Label).
			SetNillableItemValue(in.Value).
			SetNillableItemColor(in.Color).
			SetNillableItemCSS(in.Css).
			SetNillableSortOrder(sortOrderPtr).
			SetNillableDescription(in.Description).
			SetNillableState(in.State).
			Save(l.ctx)
	}
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	sortOrder := int32(dictItem.SortOrder)
	return &core.DictItemInfo{
		Id:          &dictItem.ID,
		CreatedAt:   pointer.ToInt64Ptr(dictItem.CreatedAt.UnixMilli()),
		UpdatedAt:   pointer.ToInt64Ptr(dictItem.UpdatedAt.UnixMilli()),
		Label:       &dictItem.ItemLabel,
		Value:       &dictItem.ItemValue,
		Color:       dictItem.ItemColor,
		Css:         dictItem.ItemCSS,
		SortOrder:   &sortOrder,
		Description: dictItem.Description,
		State:       &dictItem.State,
	}, nil
}

// 验证新增参数可用性
func (l *CreateOrUpdateDictItemLogic) validateCreate(in *core.DictItemInfo) error {
	if in.Label == nil || in.Value == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 验证更新参数可用性
func (l *CreateOrUpdateDictItemLogic) validateUpdate(in *core.DictItemInfo) error {
	if in.Id == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}
