package menu

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除菜单
func NewDeleteMenuLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeleteMenuLogic) DeleteMenu(req *types.DeleteMenuRequest) (resp *types.BaseResponse, err error) {
	_, err = l.svcCtx.MenuRpc.DeleteMenu(l.ctx, &core.ID32Request{Id: req.ID})
	if err != nil {
		return nil, err
	}
	resp = &types.BaseResponse{
		Code:    0,
		Message: "success",
	}
	return
}
