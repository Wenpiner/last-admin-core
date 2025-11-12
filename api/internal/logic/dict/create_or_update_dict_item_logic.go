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

type CreateOrUpdateDictItemLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建或更新字典子项
func NewCreateOrUpdateDictItemLogic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdateDictItemLogic {
	return &CreateOrUpdateDictItemLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CreateOrUpdateDictItemLogic) CreateOrUpdateDictItem(req *types.DictItemInfo) (resp *types.DictItemResponse, err error) {
	dictItemResult, err := l.svcCtx.DictRpc.CreateOrUpdateDictItem(l.ctx, &core.DictItemInfo{
		Id:          req.ID,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
		Label:       &req.Label,
		Value:       &req.Value,
		Color:       &req.Color,
		Css:         &req.Css,
		SortOrder:   &req.SortOrder,
		Description: &req.Description,
		State:       req.State,
		DictTypeId:  &req.DictID,
	})

	if err != nil {
		return nil, err
	}

	resp = &types.DictItemResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.DictItemInfo{
		ID:          dictItemResult.Id,
		CreatedAt:   dictItemResult.CreatedAt,
		UpdatedAt:   dictItemResult.UpdatedAt,
		Label:       pointer.GetString(dictItemResult.Label),
		Value:       pointer.GetString(dictItemResult.Value),
		Color:       pointer.GetString(dictItemResult.Color),
		Css:         pointer.GetString(dictItemResult.Css),
		SortOrder:   pointer.GetInt32(dictItemResult.SortOrder),
		Description: pointer.GetString(dictItemResult.Description),
		State:       dictItemResult.State,
		},
	}

	return
}
