package public_user

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wenpiner/last-admin-common/utils/encrypt"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/tokenservice"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"
	"google.golang.org/grpc/status"

	"net/http"

	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	lastHttp "github.com/wenpiner/last-admin-common/utils/http"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByPasswordLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 账号密码登录
func NewLoginByPasswordLogic(r *http.Request, svcCtx *svc.ServiceContext) *LoginByPasswordLogic {
	return &LoginByPasswordLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *LoginByPasswordLogic) LoginByPassword(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 验证验证码
	if l.svcCtx.CaptchaService.Verify(req.Captcha.ID, req.Captcha.Value, false) == false {
		return nil, errorx.NewApiInvalidParamsError("captcha.verifyFailed")
	}

	// 获取用户信息
	user, err := l.svcCtx.UserRpc.GetUserByUsername(l.ctx, &userservice.StringRequest{Value: req.Username})
	if err != nil {
		if e,ok := status.FromError(err); ok {
			if e.Message() == last_i18n.TargetNotExist {
				return nil, errorx.NewApiError(errorx.CodeInvalidCredentials, "login.passwordError")
			}
		}
		return nil, err
	}

	// 校验密码
	if encrypt.BcryptCheck(req.Password, *user.PasswordHash) == false {
		return nil, errorx.NewApiError(errorx.CodeInvalidCredentials, "login.passwordError")
	}

	// 验证用户状态
	if pointer.GetBool(user.State) == false {
		return nil, errorx.NewApiError(errorx.CodeAccountDisabled, "user.disabled")
	}

	if user.TotpInfo != nil {
		if pointer.GetBool(user.TotpInfo.State) == true {
			if req.TotpCode == nil || *req.TotpCode == "" {
				return nil, errorx.NewApiError(errorx.CodeTOTPRequired, "totp.notProvided")
			}else{
				// 验证TOTP
				verifyReq := &userservice.VerifyTotpCodeRequest{
					UserId:   *user.Id,
					TotpCode: *req.TotpCode,
				}
				verifyResp, err := l.svcCtx.UserRpc.VerifyTotpCode(l.ctx, verifyReq)
				if err != nil {
					return nil, err
				}
				if verifyResp.IsValid == false {
					return nil, errorx.NewApiError(errorx.CodeTOTPVerifyFailed, "totp.verifyFailed")
				}
			}
		}
	}

	// 生成Token
	accessToken, err := generateToken(user, l.svcCtx.Config.Auth.AccessExpire, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.TokenRpc.CreateToken(l.ctx, &tokenservice.CreateTokenRequest{
		TokenValue: accessToken,
		TokenType:  "access_token",
		UserId:     user.Id,
		ExpiresAt:  time.Now().Add(time.Second * time.Duration(l.svcCtx.Config.Auth.AccessExpire)).Unix(),
		DeviceInfo: pointer.ToStringPtr(l.r.UserAgent()),
		IpAddress:  pointer.ToStringPtr(lastHttp.GetIP(l.r)),
		UserAgent:  pointer.ToStringPtr(l.r.UserAgent()),
	})
	if err != nil {
		return nil, err
	}

	resp = &types.LoginResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "login.success",
		},
		Data: types.LoginInfo{
			AccessToken: accessToken,
		},
	}
	return
}

func generateToken(user *userservice.UserInfo, expire int64, secret string) (accessToken string, err error) {
	// 生成Token
	claims := make(jwt.MapClaims)
	iat := time.Now().Unix()
	claims["iat"] = iat
	claims["exp"] = iat + expire
	claims["userId"] = *user.Id
	claims["deptId"] = *user.DepartmentId
	claims["roleId"] = strings.Join(user.RoleValues, ",")
	claims["providerId"] = 0
	if user.ProviderId != nil {
		claims["providerId"] = *user.ProviderId
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(secret))
	if err != nil {
		err = errorx.NewInternalError("token.generateTokenFailed")
	}
	return
}
