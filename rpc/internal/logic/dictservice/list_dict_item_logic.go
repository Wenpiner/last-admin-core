package dictservicelogic

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent/dictitem"
	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictItemLogic {
	return &ListDictItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取字典子项列表
func (l *ListDictItemLogic) ListDictItem(in *core.DictItemListRequest) (*core.DictItemListResponse, error) {
	var predicates []predicate.DictItem
	if in.DictId != nil {
		predicates = append(predicates, dictitem.DictTypeIDEQ(uint32(*in.DictId)))
	}
	if in.Label != nil {
		predicates = append(predicates, dictitem.ItemLabelContains(*in.Label))
	}
	if in.Value != nil {
		predicates = append(predicates, dictitem.ItemValueContains(*in.Value))
	}

	page, err := l.svcCtx.DBEnt.DictItem.Query().Where(predicates...).Order(dictitem.BySortOrder(sql.OrderDesc())).Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	resp := &core.DictItemListResponse{
		Page: &core.BasePageResp{
			Total:      page.PageDetails.Total,
			PageNumber: page.PageDetails.Page,
			PageSize:   page.PageDetails.Pages,
		},
	}
	for _, v := range page.List {
		sortOrder := int32(v.SortOrder)
		resp.List = append(resp.List, &core.DictItemInfo{
			Id:          &v.ID,
			CreatedAt:   pointer.ToInt64Ptr(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointer.ToInt64Ptr(v.UpdatedAt.UnixMilli()),
			Label:       &v.ItemLabel,
			Value:       &v.ItemValue,
			Color:       v.ItemColor,
			Css:         v.ItemCSS,
			SortOrder:   &sortOrder,
			Description: v.Description,
			State:       &v.State,
		})
	}
	return resp, nil
}
