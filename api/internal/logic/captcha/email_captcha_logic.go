package captcha

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type EmailCaptchaLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 邮箱验证码
func NewEmailCaptchaLogic(r *http.Request, svcCtx *svc.ServiceContext) *EmailCaptchaLogic {
	return &EmailCaptchaLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *EmailCaptchaLogic) EmailCaptcha(req *types.SendEmailCaptchaReq) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
