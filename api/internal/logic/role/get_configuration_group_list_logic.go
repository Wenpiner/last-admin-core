package role

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigurationGroupListLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前系统中的所有分组列表
func NewGetConfigurationGroupListLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetConfigurationGroupListLogic {
	return &GetConfigurationGroupListLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetConfigurationGroupListLogic) GetConfigurationGroupList() (resp *types.RoleConfigurationGroupListResponse, err error) {
	response, err := l.svcCtx.RoleRpc.GetConfigurationGroup(l.ctx, &core.StringRequest{
		Value: "",
	})
	if err != nil {
		return nil, err
	}

	resp = &types.RoleConfigurationGroupListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: response.List,
	}

	return
}
