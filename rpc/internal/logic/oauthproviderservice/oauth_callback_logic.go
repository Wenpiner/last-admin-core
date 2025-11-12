package oauthproviderservicelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/user"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"
	"golang.org/x/oauth2"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type OauthCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOauthCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OauthCallbackLogic {
	return &OauthCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Oauth Callback
func (l *OauthCallbackLogic) OauthCallback(in *core.OauthCallbackRequest) (*core.UserInfo, error) {
	// 1. 验证 JWT state 参数
	claims, err := l.validateJWTState(in.State)
	if err != nil {
		return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 2. 获取 provider 配置
	providerID, ok := claims["provider_id"].(float64)
	if !ok {
		return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	providerIntID := uint32(providerID)
	provider, err := l.svcCtx.DBEnt.OauthProvider.Get(l.ctx, providerIntID)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 3. 创建 oauth2.Config
	config := &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.AuthorizationURL,
			TokenURL: provider.TokenURL,
		},
	}

	// 4. 使用 code 交换访问令牌
	token, err := config.Exchange(l.ctx, in.Code)
	if err != nil {
		l.Logger.Errorf("Failed to exchange token: %v", err)
		return nil, errorx.NewInternalError(last_i18n.DatabaseError)
	}

	// 5. 获取第三方用户信息
	oauthUserInfo, err := l.fetchUserInfo(provider, token.AccessToken)
	if err != nil {
		return nil, err
	}

	// 6. 根据邮箱和手机号查询本地用户
	localUser, err := l.findLocalUser(oauthUserInfo)
	if err != nil {
		return nil, err
	}

	// 7. 返回本地用户信息
	u := l.convertUserToUserInfo(localUser)
	u.ProviderId = &provider.ID
	return u, nil
}

// findLocalUser 根据第三方用户信息查找本地用户
func (l *OauthCallbackLogic) findLocalUser(oauthUserInfo *OAuthUserInfo) (*ent.User, error) {
	// 优先通过邮箱查找用户
	if oauthUserInfo.Email != "" {
		existingUser, err := l.svcCtx.DBEnt.User.Query().
			Where(user.EmailEQ(oauthUserInfo.Email)).
			First(l.ctx)
		if err == nil {
			return existingUser, nil
		}
		// 如果不是 NotFound 错误，返回错误
		if !ent.IsNotFound(err) {
			return nil, errorhandler.DBEntError(l.Logger, err, oauthUserInfo)
		}
	}

	// 通过手机号查找用户
	if oauthUserInfo.Phone != "" {
		existingUser, err := l.svcCtx.DBEnt.User.Query().
			Where(user.MobileEQ(oauthUserInfo.Phone)).
			First(l.ctx)
		if err == nil {
			return existingUser, nil
		}
		// 如果不是 NotFound 错误，返回错误
		if !ent.IsNotFound(err) {
			return nil, errorhandler.DBEntError(l.Logger, err, oauthUserInfo)
		}
	}

	// 如果没有找到用户，返回错误（不自动创建）
	return nil, errorx.NewInvalidArgumentError(last_i18n.TargetNotExist)
}

// convertUserToUserInfo 将用户实体转换为 UserInfo
func (l *OauthCallbackLogic) convertUserToUserInfo(userEntity *ent.User) *core.UserInfo {
	return &core.UserInfo{
		Id:           pointer.ToStringPtrIfNotEmpty(userEntity.ID.String()),
		CreatedAt:    pointer.ToInt64Ptr(userEntity.CreatedAt.UnixMilli()),
		UpdatedAt:    pointer.ToInt64Ptr(userEntity.UpdatedAt.UnixMilli()),
		Username:     &userEntity.Username,
		Email:        pointer.ToStringPtrIfNotEmpty(userEntity.Email),
		Mobile:       pointer.ToStringPtrIfNotEmpty(userEntity.Mobile),
		Avatar:       pointer.ToStringPtrIfNotEmpty(userEntity.Avatar),
		DepartmentId: &userEntity.DepartmentID,
	}
}

// OAuthUserInfo 第三方用户信息结构
type OAuthUserInfo struct {
	ID       string `json:"id"`
	Username string `json:"login,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar_url,omitempty"`
}

// validateJWTState 验证 JWT state 参数
func (l *OauthCallbackLogic) validateJWTState(stateToken string) (jwt.MapClaims, error) {
	// 使用配置的密钥
	secretKey := []byte(l.svcCtx.Config.OAuthStateSecret)

	// 解析 JWT token
	token, err := jwt.Parse(stateToken, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证 token 有效性
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 检查过期时间
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().UnixMilli() > int64(exp) {
				return nil, fmt.Errorf("token expired")
			}
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// fetchUserInfo 从第三方获取用户信息
func (l *OauthCallbackLogic) fetchUserInfo(provider *ent.OauthProvider, accessToken string) (*OAuthUserInfo, error) {
	// 如果没有配置 userinfo URL，返回空用户信息
	if provider.UserinfoURL == nil || *provider.UserinfoURL == "" {
		return &OAuthUserInfo{}, nil
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(l.ctx, "GET", *provider.UserinfoURL, nil)
	if err != nil {
		return nil, errorx.NewInternalError(last_i18n.DatabaseError)
	}

	// 设置授权头
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		l.Logger.Errorf("Failed to fetch user info: %v", err)
		return nil, errorx.NewInternalError(last_i18n.DatabaseError)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errorx.NewInternalError(last_i18n.DatabaseError)
	}

	// 解析 JSON
	var userInfo OAuthUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		l.Logger.Errorf("Failed to parse user info: %v, body: %s", err, string(body))
		return nil, errorx.NewInternalError(last_i18n.DatabaseError)
	}

	return &userInfo, nil
}
