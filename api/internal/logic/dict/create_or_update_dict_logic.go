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

type CreateOrUpdateDictLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建或更新字典
func NewCreateOrUpdateDictLogic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdateDictLogic {
	return &CreateOrUpdateDictLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CreateOrUpdateDictLogic) CreateOrUpdateDict(req *types.DictInfo) (resp *types.DictInfoResponse, err error) {
	dictResult, err := l.svcCtx.DictRpc.CreateOrUpdateDict(l.ctx, &core.DictInfo{
		Id:          req.ID,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
		Name:        &req.Name,
		Code:        &req.Code,
		Description: &req.Description,
		State:       req.State,
	})

	if err != nil {
		return nil, err
	}

	resp = &types.DictInfoResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.DictInfo{
			ID:          dictResult.Id,
			CreatedAt:   dictResult.CreatedAt,
			UpdatedAt:   dictResult.UpdatedAt,
			Name:        pointer.GetString(dictResult.Name),
			Code:        pointer.GetString(dictResult.Code),
			Description: pointer.GetString(dictResult.Description),
			State:       dictResult.State,
		},
	}

	return
}
