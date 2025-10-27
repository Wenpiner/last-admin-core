package tokenservicelogic

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/token"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateTokenLogic {
	return &ValidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 验证Token
func (l *ValidateTokenLogic) ValidateToken(in *core.ValidateTokenRequest) (*core.ValidateTokenResponse, error) {
	// 查询Token
	query := l.svcCtx.DBEnt.Token.Query().
		Where(token.TokenValueEQ(in.TokenValue))

	// 如果指定了token类型，添加类型过滤
	if in.TokenType != nil && *in.TokenType != "" {
		query = query.Where(token.TokenTypeEQ(*in.TokenType))
	}

	tokenEntity, err := query.First(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return &core.ValidateTokenResponse{
				IsValid: false,
				Message: pointer.ToStringPtr("Token not found"),
			}, nil
		}
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 检查Token是否已撤销
	if tokenEntity.IsRevoked {
		return &core.ValidateTokenResponse{
			IsValid: false,
			Message: pointer.ToStringPtr("Token has been revoked"),
		}, nil
	}

	// 检查Token是否过期
	if time.Now().After(tokenEntity.ExpiresAt) {
		return &core.ValidateTokenResponse{
			IsValid: false,
			Message: pointer.ToStringPtr("Token has expired"),
		}, nil
	}

	// Token有效
	return &core.ValidateTokenResponse{
		IsValid:   true,
		TokenInfo: l.convertTokenToTokenInfo(tokenEntity),
		Message:   pointer.ToStringPtr("Token is valid"),
	}, nil
}

// convertTokenToTokenInfo 将Token实体转换为TokenInfo
func (l *ValidateTokenLogic) convertTokenToTokenInfo(tokenEntity *ent.Token) *core.TokenInfo {
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
func (l *ValidateTokenLogic) convertUserIDToString(userID *uuid.UUID) *string {
	if userID == nil {
		return nil
	}
	userIDStr := userID.String()
	return &userIDStr
}

// convertTimeToInt64Ptr 将时间转换为int64指针
func (l *ValidateTokenLogic) convertTimeToInt64Ptr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	timestamp := t.UnixMilli()
	return &timestamp
}
