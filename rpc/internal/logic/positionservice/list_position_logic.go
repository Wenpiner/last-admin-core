package positionservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/position"
	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPositionLogic {
	return &ListPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取岗位列表
func (l *ListPositionLogic) ListPosition(in *core.PositionListRequest) (*core.PositionListResponse, error) {
	// 构建查询条件
	var predicates []predicate.Position

	// 根据岗位名称模糊搜索
	if in.PositionName != nil && *in.PositionName != "" {
		predicates = append(predicates, position.PositionNameContains(*in.PositionName))
	}

	// 根据岗位编码模糊搜索
	if in.PositionCode != nil && *in.PositionCode != "" {
		predicates = append(predicates, position.PositionCodeContains(*in.PositionCode))
	}

	page, err := l.svcCtx.DBEnt.Position.Query().Where(predicates...).Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	resp := &core.PositionListResponse{
		Page: &core.BasePageResp{
			Total:      page.PageDetails.Total,
			PageNumber: in.Page.PageNumber,
			PageSize:   in.Page.PageSize,
		},
	}

	for _, pos := range page.List {
		resp.List = append(resp.List, l.convertPositionToPositionInfo(pos))
	}

	return resp, nil
}

// 将 Position 实体转换为 PositionInfo
func (l *ListPositionLogic) convertPositionToPositionInfo(pos *ent.Position) *core.PositionInfo {
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
