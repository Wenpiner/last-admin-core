package auth

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccessCodesLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户权限码(通过Menu获取按钮级别的权限)
func NewGetAccessCodesLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetAccessCodesLogic {
	return &GetAccessCodesLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetAccessCodesLogic) GetAccessCodes() (resp *types.AccessCodesResponse, err error) {
	res, err := l.svcCtx.MenuRpc.ListPagePermissionByRole(l.ctx, &core.StringRequest{Value: l.ctx.Value("roleId").(string)})
	if err != nil {
		return nil, err
	}

	resp = &types.AccessCodesResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: res.List,
	}

	return
}
