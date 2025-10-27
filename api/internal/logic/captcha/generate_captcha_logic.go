package captcha

import (
	"context"
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

// 生成验证码
func NewGenerateCaptchaLogic(r *http.Request, svcCtx *svc.ServiceContext) *GenerateCaptchaLogic {
	return &GenerateCaptchaLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		ctx:    r.Context(),
		svcCtx: svcCtx,
	}
}

func (l *GenerateCaptchaLogic) GenerateCaptcha() (resp *types.GenerateCaptchaResp, err error) {
	result, err := l.svcCtx.CaptchaService.Generate()
	if err != nil {
		return nil, errorx.NewInternalError("captcha.generateCaptchaFailed")
	}

	resp = &types.GenerateCaptchaResp{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.GenerateCaptchaInfo{
			ID:         result.ID,
			Base64Blob: result.Base64Blob,
			CaptchaType: result.CaptchaType,
		},
	}

	return
}
