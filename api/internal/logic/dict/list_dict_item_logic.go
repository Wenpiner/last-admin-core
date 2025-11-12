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

type ListDictItemLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取字典子项列表
func NewListDictItemLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListDictItemLogic {
	return &ListDictItemLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListDictItemLogic) ListDictItem(req *types.DictItemListRequest) (resp *types.DictItemListResponse, err error) {
	dictItemList, err := l.svcCtx.DictRpc.ListDictItem(l.ctx, &core.DictItemListRequest{
		Page: &core.BasePageRequest{
			PageNumber: req.Page.CurrentPage,
			PageSize:   req.Page.PageSize,
		},
		DictId: &req.DictId,
		Label:  &req.Label,
		Value:  &req.Value,
	})

	if err != nil {
		return nil, err
	}

	resp = &types.DictItemListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.DictItemListInfo{
			BaseListInfo: types.BaseListInfo{
				Total: dictItemList.Page.Total,
			},
			List: ConvertToDictItemInfo(dictItemList.List),
		},
	}

	return
}

func ConvertToDictItemInfo(dictItemList []*core.DictItemInfo) []types.DictItemInfo {
	var result []types.DictItemInfo
	for _, item := range dictItemList {
		result = append(result, types.DictItemInfo{
			ID:          item.Id,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			Label:       pointer.GetString(item.Label),
			Value:       pointer.GetString(item.Value),
			Color:       pointer.GetString(item.Color),
			Css:         pointer.GetString(item.Css),
			SortOrder:   pointer.GetInt32(item.SortOrder),
			Description: pointer.GetString(item.Description),
			State:       item.State,
		})
	}
	return result
}
