package dict

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取字典列表
func NewListDictLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListDictLogic {
	return &ListDictLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListDictLogic) ListDict(req *types.DictListRequest) (resp *types.DictListResponse, err error) {
	dictList, err := l.svcCtx.DictRpc.ListDict(l.ctx, &core.DictListRequest{
		Page: &core.BasePageRequest{
			PageNumber: req.Page.CurrentPage,
			PageSize:   req.Page.PageSize,
		},
		Name:        &req.Name,
		Code:        &req.Code,
		Description: &req.Description,
	})

	if err != nil {
		return nil, err
	}

	resp = &types.DictListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.DictListInfo{
			BaseListInfo: types.BaseListInfo{
				Total: dictList.Page.Total,
			},
			List: ConvertToDictInfo(dictList.List),
		},
	}

	return
}

func ConvertToDictInfo(dictList []*core.DictInfo) []types.DictInfo {
	var result []types.DictInfo
	for _, dict := range dictList {
		result = append(result, types.DictInfo{
			ID:          dict.Id,
			CreatedAt:   dict.CreatedAt,
			UpdatedAt:   dict.UpdatedAt,
			Name:        pointer.GetString(dict.Name),
			Code:        pointer.GetString(dict.Code),
			Description: pointer.GetString(dict.Description),
			State:       dict.State,
		})
	}
	return result
}
