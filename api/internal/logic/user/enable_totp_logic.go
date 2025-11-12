package user

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnableTotpLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建/重置TOTP
func NewEnableTotpLogic(r *http.Request, svcCtx *svc.ServiceContext) *EnableTotpLogic {
	return &EnableTotpLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *EnableTotpLogic) EnableTotp(req *types.EnableTotpRequest) (resp *types.TotpSetupResponse, err error) {
	response, err := l.svcCtx.UserRpc.EnableTotp(l.ctx, &userservice.EnableTotpRequest{
		UserId: req.UserId,
		Domain: req.Domain,
		Issuer: req.Issuer,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.TotpSetupResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "totp.setupSuccess",
		},
		Data: types.TotpSetupInfo{
			QRText: response.QrCodeContent,
		},
	}

	return
}
