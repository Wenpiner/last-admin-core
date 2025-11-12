package api

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateApiLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建或更新API
func NewCreateOrUpdateApiLogic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdateApiLogic {
	return &CreateOrUpdateApiLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CreateOrUpdateApiLogic) CreateOrUpdateApi(req *types.ApiInfo) (resp *types.ApiInfo, err error) {
	apiResult, err := l.svcCtx.ApiRpc.CreateOrUpdateApi(l.ctx, &core.ApiInfo{
		Id:          &req.ID,
		CreatedAt:   &req.CreatedAt,
		UpdatedAt:   &req.UpdatedAt,
		Name:        &req.Name,
		Method:      &req.Method,
		Path:        &req.Path,
		Description: &req.Description,
		IsRequired:  &req.IsRequired,
		ServiceName: &req.ServiceName,
		ApiGroup:    &req.ApiGroup,
	})

	if err != nil {
		return nil, err
	}

	resp = &types.ApiInfo{
		ID:          pointer.GetUint32(apiResult.Id),
		CreatedAt:   *apiResult.CreatedAt,
		UpdatedAt:   *apiResult.UpdatedAt,
		Name:        pointer.GetString(apiResult.Name),
		Method:      pointer.GetString(apiResult.Method),
		Path:        pointer.GetString(apiResult.Path),
		Description: pointer.GetString(apiResult.Description),
		IsRequired:  *apiResult.IsRequired,
		ServiceName: pointer.GetString(apiResult.ServiceName),
		ApiGroup:    pointer.GetString(apiResult.ApiGroup),
	}

	return
}
