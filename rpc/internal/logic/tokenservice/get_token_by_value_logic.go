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

type GetTokenByValueLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenByValueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByValueLogic {
	return &GetTokenByValueLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据Token值获取Token信息
func (l *GetTokenByValueLogic) GetTokenByValue(in *core.StringRequest) (*core.TokenInfo, error) {
	// 查询Token
	tokenEntity, err := l.svcCtx.DBEnt.Token.Query().
		Where(token.TokenValueEQ(in.Value)).
		First(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertTokenToTokenInfo(tokenEntity), nil
}

// convertTokenToTokenInfo 将Token实体转换为TokenInfo
func (l *GetTokenByValueLogic) convertTokenToTokenInfo(tokenEntity *ent.Token) *core.TokenInfo {
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
func (l *GetTokenByValueLogic) convertUserIDToString(userID *uuid.UUID) *string {
	if userID == nil {
		return nil
	}
	userIDStr := userID.String()
	return &userIDStr
}

// convertTimeToInt64Ptr 将时间转换为int64指针
func (l *GetTokenByValueLogic) convertTimeToInt64Ptr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	timestamp := t.UnixMilli()
	return &timestamp
}
