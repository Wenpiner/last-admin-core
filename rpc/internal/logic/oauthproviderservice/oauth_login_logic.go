package oauthproviderservicelogic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent/oauthprovider"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"
	"golang.org/x/oauth2"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type OauthLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOauthLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OauthLoginLogic {
	return &OauthLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Oauth Login
func (l *OauthLoginLogic) OauthLogin(in *core.OauthLoginRequest) (*core.OauthRedirectResponse, error) {
	// 1. 根据 provider 查询 OAuth 配置
	provider, err := l.svcCtx.DBEnt.OauthProvider.Query().
		Where(oauthprovider.ProviderCodeEQ(in.Provider)).
		Where(oauthprovider.StateEQ(true)).
		Only(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 2. 创建 oauth2.Config
	config := &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURI,
		Scopes:       l.parseScopes(provider.Scopes),
		Endpoint: oauth2.Endpoint{
			AuthURL:   provider.AuthorizationURL,
			TokenURL:  provider.TokenURL,
			AuthStyle: oauth2.AuthStyle(provider.AuthStyle),
		},
	}

	// 3. 生成 JWT state 并包含必要信息
	state, err := l.generateJWTState(in.State, provider.ID)
	if err != nil {
		return nil, errorx.NewInternalError(last_i18n.DatabaseError)
	}

	// 4. 生成授权URL
	authURL := config.AuthCodeURL(state)

	return &core.OauthRedirectResponse{
		Url: authURL,
	}, nil
}

// parseScopes 解析 scopes 字符串为数组
func (l *OauthLoginLogic) parseScopes(scopes *string) []string {
	if scopes == nil || *scopes == "" {
		return []string{}
	}
	// 使用空格分割 scopes
	return strings.Split(*scopes, " ")
}

// generateJWTState 生成包含状态信息的 JWT token
func (l *OauthLoginLogic) generateJWTState(originalState string, providerID uint32) (string, error) {
	// 生成随机字符串作为额外的安全措施
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(randomBytes)

	// 创建 JWT claims
	claims := jwt.MapClaims{
		"state":       originalState,
		"provider_id": providerID,
		"nonce":       nonce,
		"exp":         time.Now().Add(10 * time.Minute).Unix(), // 10分钟过期
		"iat":         time.Now().Unix(),
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用配置的密钥签名
	secretKey := []byte(l.svcCtx.Config.OAuthStateSecret)

	return token.SignedString(secretKey)
}
