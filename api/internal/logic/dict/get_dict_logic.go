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

type GetDictLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取字典
func NewGetDictLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetDictLogic {
	return &GetDictLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetDictLogic) GetDict(req *types.ID32Request) (resp *types.DictInfoResponse, err error) {
	dictResult, err := l.svcCtx.DictRpc.GetDict(l.ctx, &core.ID32Request{Id: req.ID})
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
