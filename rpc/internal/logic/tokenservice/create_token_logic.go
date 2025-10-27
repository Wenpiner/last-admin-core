package tokenservicelogic

import (
	"context"
	"time"

	"github.com/google/uuid"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTokenLogic {
	return &CreateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建Token
func (l *CreateTokenLogic) CreateToken(in *core.CreateTokenRequest) (*core.TokenInfo, error) {
	// 验证必填字段
	if err := l.validateCreateTokenRequest(in); err != nil {
		return nil, err
	}

	// 解析过期时间
	expiresAt := time.Unix(in.ExpiresAt, 0)

	// 解析用户ID
	var userID *uuid.UUID
	if in.UserId != nil && *in.UserId != "" {
		parsedUserID, err := uuid.Parse(*in.UserId)
		if err != nil {
			return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		userID = &parsedUserID
	}

	// 创建Token
	tokenEntity, err := l.svcCtx.DBEnt.Token.Create().
		SetTokenValue(in.TokenValue).
		SetTokenType(in.TokenType).
		SetNillableUserID(userID).
		SetExpiresAt(expiresAt).
		SetNillableDeviceInfo(in.DeviceInfo).
		SetNillableIPAddress(in.IpAddress).
		SetNillableUserAgent(in.UserAgent).
		SetNillableMetadata(in.Metadata).
		SetNillableRefreshTokenID(in.RefreshTokenId).
		Save(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertTokenToTokenInfo(tokenEntity), nil
}

// validateCreateTokenRequest 验证创建Token请求
func (l *CreateTokenLogic) validateCreateTokenRequest(in *core.CreateTokenRequest) error {
	if in.TokenValue == "" {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	if in.TokenType == "" {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	if in.ExpiresAt <= 0 {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// convertTokenToTokenInfo 将Token实体转换为TokenInfo
func (l *CreateTokenLogic) convertTokenToTokenInfo(tokenEntity *ent.Token) *core.TokenInfo {
	return &core.TokenInfo{
		Id:             pointer.ToUint32Ptr(uint32(tokenEntity.ID)),
		CreatedAt:      pointer.ToInt64Ptr(tokenEntity.CreatedAt.UnixMilli()),
		UpdatedAt:      pointer.ToInt64Ptr(tokenEntity.UpdatedAt.UnixMilli()),
		State:          &tokenEntity.State,
		TokenValue:     &tokenEntity.TokenValue,
		TokenType:      &tokenEntity.TokenType,
		UserId:         l.convertUserIDToString(tokenEntity.UserID),
		ExpiresAt:      pointer.ToInt64Ptr(tokenEntity.ExpiresAt.UnixMilli()),
		IsRevoked:      &tokenEntity.IsRevoked,
		DeviceInfo:     tokenEntity.DeviceInfo,
		IpAddress:      tokenEntity.IPAddress,
		LastUsedAt:     l.convertTimeToInt64Ptr(tokenEntity.LastUsedAt),
		UserAgent:      tokenEntity.UserAgent,
		Metadata:       tokenEntity.Metadata,
		RefreshTokenId: tokenEntity.RefreshTokenID,
	}
}

// convertUserIDToString 将UUID转换为字符串指针
func (l *CreateTokenLogic) convertUserIDToString(userID *uuid.UUID) *string {
	if userID == nil {
		return nil
	}
	userIDStr := userID.String()
	return &userIDStr
}

// convertTimeToInt64Ptr 将时间转换为int64指针
func (l *CreateTokenLogic) convertTimeToInt64Ptr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	timestamp := t.UnixMilli()
	return &timestamp
}
