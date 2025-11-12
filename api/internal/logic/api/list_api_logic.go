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

type ListApiLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取API列表
func NewListApiLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListApiLogic {
	return &ListApiLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListApiLogic) ListApi(req *types.ApiListRequest) (resp *types.ApiListResponse, err error) {
	apiList, err := l.svcCtx.ApiRpc.ListApi(l.ctx, &core.ApiListRequest{
		Page: &core.BasePageRequest{
			PageNumber: req.Page.CurrentPage,
			PageSize:   req.Page.PageSize,
		},
		ServiceName: &req.ServiceName,
		ApiGroup:    &req.ApiGroup,
		Method:      &req.Method,
		Description: &req.Description,
		Path:        &req.Path,
	})

	if err != nil {
		return nil, err
	}

	resp = &types.ApiListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.ApiListInfo{
			BaseListInfo: types.BaseListInfo{
				Total: apiList.Page.Total,
			},
			List: ConvertToApiInfo(apiList.List),
		},
	}

	return
}

func ConvertToApiInfo(apiList []*core.ApiInfo) []types.ApiInfo {
	var result []types.ApiInfo
	for _, api := range apiList {
		result = append(result, types.ApiInfo{
			ID:          pointer.GetUint32(api.Id),
			CreatedAt:   *api.CreatedAt,
			UpdatedAt:   *api.UpdatedAt,
			Name:        pointer.GetString(api.Name),
			Method:      pointer.GetString(api.Method),
			Path:        pointer.GetString(api.Path),
			Description: pointer.GetString(api.Description),
			IsRequired:  *api.IsRequired,
			ServiceName: pointer.GetString(api.ServiceName),
			ApiGroup:    pointer.GetString(api.ApiGroup),
		})
	}
	return result
}
