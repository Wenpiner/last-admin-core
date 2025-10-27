package positionservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionLogic {
	return &GetPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取岗位
func (l *GetPositionLogic) GetPosition(in *core.ID32Request) (*core.PositionInfo, error) {
	pos, err := l.svcCtx.DBEnt.Position.Get(l.ctx, in.Id)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertPositionToPositionInfo(pos), nil
}

// 将 Position 实体转换为 PositionInfo
func (l *GetPositionLogic) convertPositionToPositionInfo(pos *ent.Position) *core.PositionInfo {
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
