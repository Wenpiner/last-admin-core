package captcha

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type SmsCaptchaLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送短信验证码
func NewSmsCaptchaLogic(r *http.Request, svcCtx *svc.ServiceContext) *SmsCaptchaLogic {
	return &SmsCaptchaLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *SmsCaptchaLogic) SmsCaptcha(req *types.SendSmsCaptchaReq) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
