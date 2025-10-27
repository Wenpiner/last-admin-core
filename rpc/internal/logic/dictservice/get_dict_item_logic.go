package dictservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictItemLogic {
	return &GetDictItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取字典子项
func (l *GetDictItemLogic) GetDictItem(in *core.ID32Request) (*core.DictItemInfo, error) {
	dictItem, err := l.svcCtx.DBEnt.DictItem.Get(l.ctx, in.Id)
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
