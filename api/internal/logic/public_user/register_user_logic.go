package public_user

import (
	"context"

	"github.com/wenpiner/last-admin-common/enums"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterUserLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注册用户
func NewRegisterUserLogic(r *http.Request, svcCtx *svc.ServiceContext) *RegisterUserLogic {
	return &RegisterUserLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *RegisterUserLogic) RegisterUser(req *types.RegisterRequest) (resp *types.BaseResponse, err error) {
	// 检查是否开启注册
	if l.svcCtx.Config.ProjectConf.OpenRegister == false {
		return nil, errorx.NewInvalidArgumentError("register.registerClosed")
	}

	// 验证验证码
	if l.svcCtx.CaptchaService.VerifyAndClear(req.CaptchaInfo.ID, req.CaptchaInfo.Value) == false {
		return nil, errorx.NewInvalidArgumentError("captcha.verifyFailed")
	}

	// 创建用户
	_, err = l.svcCtx.UserRpc.CreateUser(l.ctx, &userservice.UserInfo{
		Username:     &req.Username,
		PasswordHash: &req.Password,
		RoleValues:   []string{l.svcCtx.Config.ProjectConf.RegisterRoleValue},
		PositionIds:  []uint32{enums.DefaultPositionID},
		DepartmentId: pointer.ToUint32Ptr(enums.DefaultDeptID),
	})
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "register.registerSuccess",
	}

	return
}
