package role

import (
	"context"
	"fmt"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleApiLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色API
func NewGetRoleApiLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetRoleApiLogic {
	return &GetRoleApiLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetRoleApiLogic) GetRoleApi(req *types.StringIDRequest) (resp *types.RoleApiListResponse, err error) {
	role, err := l.svcCtx.RoleRpc.GetRoleByValue(l.ctx, &core.StringRequest{
		Value: req.ID,
	})
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.Casbin.GetFilteredPolicy(0, *role.RoleCode)
	if err != nil {
		return nil, err
	}

	var apiList []string
	for _, v := range result {
		apiList = append(apiList, fmt.Sprintf("%s|%s", v[2], v[1]))
	}

	resp = &types.RoleApiListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: apiList,
	}
	return
}
