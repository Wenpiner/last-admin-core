package user

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyTotpLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 验证TOTP
func NewVerifyTotpLogic(r *http.Request, svcCtx *svc.ServiceContext) *VerifyTotpLogic {
	return &VerifyTotpLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *VerifyTotpLogic) VerifyTotp(req *types.VerifyTotpRequest) (resp *types.BaseResponse, err error) {
	response, err := l.svcCtx.UserRpc.VerifyTotpSetup(l.ctx, &userservice.VerifyTotpSetupRequest{
		UserId:   req.UserId,
		TotpCode: req.Code,
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
