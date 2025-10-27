package apiservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiLogic {
	return &GetApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取API
func (l *GetApiLogic) GetApi(in *core.ID32Request) (*core.ApiInfo, error) {
	// 查询API
	apiEntity, err := l.svcCtx.DBEnt.API.Get(l.ctx, in.Id)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertAPIToApiInfo(apiEntity), nil
}

// 将 API 实体转换为 ApiInfo
func (l *GetApiLogic) convertAPIToApiInfo(apiEntity *ent.API) *core.ApiInfo {
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
