package configuration

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConfigurationLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除配置
func NewDeleteConfigurationLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteConfigurationLogic {
	return &DeleteConfigurationLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeleteConfigurationLogic) DeleteConfiguration(req *types.DeleteConfigurationRequest) (resp *types.BaseResponse, err error) {
	_, err = l.svcCtx.ConfigurationRpc.DeleteConfiguration(l.ctx, &core.StringRequest{Value: req.Key})
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "success",
	}

	return
}
