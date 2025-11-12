package tokenservicelogic

import (
	"time"

	"github.com/google/uuid"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/types/core"
)

// ConvertTokenToTokenInfo 将Token实体转换为TokenInfo
func ConvertTokenToTokenInfo(tokenEntity *ent.Token) *core.TokenInfo {
	t := core.TokenInfo{
		Id:             pointer.ToUint32Ptr(uint32(tokenEntity.ID)),
		CreatedAt:      pointer.ToInt64Ptr(tokenEntity.CreatedAt.UnixMilli()),
		UpdatedAt:      pointer.ToInt64Ptr(tokenEntity.UpdatedAt.UnixMilli()),
		State:          &tokenEntity.State,
		TokenValue:     &tokenEntity.TokenValue,
		TokenType:      &tokenEntity.TokenType,
		UserId:         ConvertUserIDToString(tokenEntity.UserID),
		ExpiresAt:      pointer.ToInt64Ptr(tokenEntity.ExpiresAt.UnixMilli()),
		DeviceInfo:     tokenEntity.DeviceInfo,
		IpAddress:      tokenEntity.IPAddress,
		LastUsedAt:     ConvertTimeToInt64Ptr(tokenEntity.LastUsedAt),
		UserAgent:      tokenEntity.UserAgent,
		Metadata:       tokenEntity.Metadata,
		ProviderId:     tokenEntity.ProviderID,
	}
	if tokenEntity.Edges.Provider != nil {
		t.ProviderName = &tokenEntity.Edges.Provider.ProviderName
	}

	if tokenEntity.Edges.User != nil {
		t.Username = &tokenEntity.Edges.User.Username
		t.FullName = &tokenEntity.Edges.User.FullName
	}
	return &t
}

// ConvertUserIDToString 将UUID转换为字符串指针
func ConvertUserIDToString(userID *uuid.UUID) *string {
	if userID == nil {
		return nil
	}
	userIDStr := userID.String()
	return &userIDStr
}

// ConvertTimeToInt64Ptr 将时间转换为int64指针
func ConvertTimeToInt64Ptr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	timestamp := t.UnixMilli()
	return &timestamp
}

