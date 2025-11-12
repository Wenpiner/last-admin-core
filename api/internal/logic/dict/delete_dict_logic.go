package dict

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除字典
func NewDeleteDictLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteDictLogic {
	return &DeleteDictLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeleteDictLogic) DeleteDict(req *types.ID32Request) (resp *types.BaseResponse, err error) {
	_, err = l.svcCtx.DictRpc.DeleteDict(l.ctx, &core.ID32Request{Id: req.ID})
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "success",
	}

	return
}
