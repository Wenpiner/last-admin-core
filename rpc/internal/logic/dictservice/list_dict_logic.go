package dictservicelogic

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent/dicttype"
	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictLogic {
	return &ListDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取字典列表
func (l *ListDictLogic) ListDict(in *core.DictListRequest) (*core.DictListResponse, error) {
	var predicates []predicate.DictType
	if in.Name != nil {
		predicates = append(predicates, dicttype.DictTypeNameContains(*in.Name))
	}
	if in.Code != nil {
		predicates = append(predicates, dicttype.DictTypeCodeContains(*in.Code))
	}
	if in.Description != nil {
		predicates = append(predicates, dicttype.DescriptionContains(*in.Description))
	}

	page, err := l.svcCtx.DBEnt.DictType.Query().Where(predicates...).Order(dicttype.ByCreatedAt(sql.OrderDesc())).Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	resp := &core.DictListResponse{
		Page: &core.BasePageResp{
			Total:      page.PageDetails.Total,
			PageNumber: page.PageDetails.Page,
			PageSize:   page.PageDetails.Pages,
		},
	}
	for _, v := range page.List {
		resp.List = append(resp.List, &core.DictInfo{
			Id:          &v.ID,
			CreatedAt:   pointer.ToInt64Ptr(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointer.ToInt64Ptr(v.UpdatedAt.UnixMilli()),
			Name:        &v.DictTypeName,
			Code:        &v.DictTypeCode,
			Description: &v.Description,
			State:       &v.State,
		})
	}
	return resp, nil
}
