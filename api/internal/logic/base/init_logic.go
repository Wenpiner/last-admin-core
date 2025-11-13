package base

import (
	"context"
	"net/http"

	"github.com/wenpiner/last-admin-common/enums"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/initservice"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type InitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewInitLogic(r *http.Request, svcCtx *svc.ServiceContext) *InitLogic {
	return &InitLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		r:      r,
		svcCtx: svcCtx,
	}
}

func (l *InitLogic) Init() (resp *types.BaseDataInfo, err error) {
	// 校验是否开启初始化
	response, _ := l.svcCtx.ConfigurationRpc.GetConfiguration(l.ctx, &core.StringRequest{
		Value: enums.ConfigurationInit,
	})
	// 检查是否开启初始化
	if response != nil && response.Value == "true" {
		return nil, errorx.NewApiError(errorx.CodeForbidden, "init.closed")
	}

	_, err = l.svcCtx.InitRpc.Init(l.ctx, &initservice.EmptyRequest{})

	return
}
