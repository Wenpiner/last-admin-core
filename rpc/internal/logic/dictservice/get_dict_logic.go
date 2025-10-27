package dictservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictLogic {
	return &GetDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取字典
func (l *GetDictLogic) GetDict(in *core.ID32Request) (*core.DictInfo, error) {
	dictType, err := l.svcCtx.DBEnt.DictType.Get(l.ctx, in.Id)
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
