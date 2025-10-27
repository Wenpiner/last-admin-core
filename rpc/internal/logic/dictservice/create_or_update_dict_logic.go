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

type CreateOrUpdateDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdateDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdateDictLogic {
	return &CreateOrUpdateDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建或更新字典
func (l *CreateOrUpdateDictLogic) CreateOrUpdateDict(in *core.DictInfo) (*core.DictInfo, error) {
	// 开启事物，并先进行检查是否存在，如果存在则进行更新否则进行创建
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()
	var dictType *ent.DictType
	if in.Id == nil {
		// 新增,验证必填参数可用性
		if err := l.validateCreate(in); err != nil {
			return nil, err
		}
		dictType, err = tx.DictType.Create().
			SetDictTypeName(pointer.GetString(in.Name)).
			SetDictTypeCode(pointer.GetString(in.Code)).
			SetDescription(pointer.GetString(in.Description)).
			SetState(pointer.GetBool(in.State)).
			Save(l.ctx)
	} else {
		// 更新
		if err := l.validateUpdate(in); err != nil {
			return nil, err
		}
		dictType, err = tx.DictType.UpdateOneID(pointer.GetUint32(in.Id)).
			SetNillableDictTypeName(in.Name).
			SetNillableDictTypeCode(in.Code).
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

	return &core.DictInfo{
		Id:          &dictType.ID,
		CreatedAt:   pointer.ToInt64Ptr(dictType.CreatedAt.UnixMilli()),
		UpdatedAt:   pointer.ToInt64Ptr(dictType.UpdatedAt.UnixMilli()),
		Name:        &dictType.DictTypeName,
		Code:        &dictType.DictTypeCode,
		Description: &dictType.Description,
		State:       &dictType.State,
	}, nil
}

// 验证新增参数可用性
func (l *CreateOrUpdateDictLogic) validateCreate(in *core.DictInfo) error {
	if in.Name == nil || in.Code == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 验证更新参数可用性
func (l *CreateOrUpdateDictLogic) validateUpdate(in *core.DictInfo) error {
	if in.Id == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}
