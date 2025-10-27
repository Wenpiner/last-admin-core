package public_user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/errorx"
)

const (
	totpVerifyPrefix = "totp:verify:"
)

type TotpScene string

const (
	// 登录
	TotpSceneLogin TotpScene = "login"
)

// 缓存进redis的结构题
type TotpInfo struct {
	UserID string `json:"user_id"`
	// 认证场景
	Scene TotpScene `json:"scene"`
	// 认证数据
	RequestInfo string `json:"request_info"`
}

func (t *TotpInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TotpInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

// CreateTotpVerify 存在Totp，进行验证
func CreateTotpVerify(ctx context.Context, userID string, redis *redis.Client, scene TotpScene, requestInfo any) (key string, expiredAt time.Time, err error) {
	// 随机生成UUID V7
	id, _ := uuid.NewV7()
	key = totpVerifyPrefix + id.String()
	// 设置过期时间
	expiredAt = time.Now().Add(5 * time.Minute)
	// 格式化requestInfo
	requestInfoJSON, err := json.Marshal(requestInfo)
	if err != nil {
		return "", time.Time{}, errorx.NewInternalError("totp.verifyFailed")
	}
	// 保存到Redis
	totpInfo := &TotpInfo{
		UserID:      userID,
		Scene:       scene,
		RequestInfo: string(requestInfoJSON),
	}
	err = redis.SetEx(ctx, key, totpInfo, expiredAt.Sub(time.Now())).Err()
	if err != nil {
		return "", time.Time{}, errorx.NewInternalError("totp.verifyFailed")
	}

	return id.String(), expiredAt, nil
}

// 通过ID获取用户认证的情况
func GetTotpByID(ctx context.Context, id string, redis *redis.Client) (*TotpInfo, error) {
	key := totpVerifyPrefix + id
	// 判断是否已过期
	expire, err := redis.TTL(ctx, key).Result()
	if err != nil {
		return nil, errorx.NewInternalError("totp.verifyFailed")
	}
	if expire < 0 {
		return nil, errorx.NewInvalidArgumentError("totp.verifyExpired")
	}

	// 获取数据
	data, err := redis.Get(ctx, key).Result()
	if err != nil {
		return nil, errorx.NewInternalError("totp.verifyFailed")
	}

	// 解析数据
	var totpInfo TotpInfo
	err = totpInfo.UnmarshalBinary([]byte(data))
	if err != nil {
		return nil, errorx.NewInternalError("totp.verifyFailed")
	}

	return &totpInfo, nil
}
