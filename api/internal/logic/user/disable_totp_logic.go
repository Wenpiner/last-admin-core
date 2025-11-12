package user

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisableTotpLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 禁用TOTP
func NewDisableTotpLogic(r *http.Request, svcCtx *svc.ServiceContext) *DisableTotpLogic {
	return &DisableTotpLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DisableTotpLogic) DisableTotp(req *types.UUIDRequest) (resp *types.BaseResponse, err error) {
	response, err := l.svcCtx.UserRpc.DisableTotp(l.ctx, &userservice.DisableTotpRequest{
		UserId: req.ID,
	})
	if err != nil {
		return nil, err
	}


	resp = &types.BaseResponse{
		Code:    0,
		Message: response.Message,
	}

	return
}
