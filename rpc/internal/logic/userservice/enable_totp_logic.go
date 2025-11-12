package userservicelogic

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent/user"
	"github.com/wenpiner/last-admin-core/rpc/ent/usertotp"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type EnableTotpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEnableTotpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnableTotpLogic {
	return &EnableTotpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// TOTP相关接口
func (l *EnableTotpLogic) EnableTotp(in *core.EnableTotpRequest) (*core.TotpSetupResponse, error) {
	// 解析用户ID
	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 检查用户是否存在
	userExists, err := l.svcCtx.DBEnt.User.Query().
		Where(user.IDEQ(userID)).
		Exist(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	if !userExists {
		return nil, errorx.NewInvalidArgumentError(last_i18n.TargetNotExist)
	}

	// 检查用户是否已经启用了TOTP
	existingTotp, err := l.svcCtx.DBEnt.UserTotp.Query().
		Where(usertotp.IDEQ(userID)).
		First(l.ctx)
	if err == nil && existingTotp.State {
		return nil, errorx.NewInvalidArgumentError("totp.alreadyEnabled")
	}

	// 验证域和发行者是否有效
	if in.Domain == "" || in.Issuer == "" {
		return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 获取用户信息用于生成TOTP账户名
	userInfo, err := l.svcCtx.DBEnt.User.Get(l.ctx, userID)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 生成TOTP密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      in.Issuer,
		AccountName: fmt.Sprintf("%s@%s", userInfo.Username, in.Domain),
	})
	if err != nil {
		return nil, errorx.NewInternalError("totp.generateSecretFailed")
	}

	// 生成备用恢复码
	backupCodes := l.generateBackupCodes(8)
	backupCodesJSON, err := json.Marshal(backupCodes)
	if err != nil {
		return nil, errorx.NewInternalError("totp.generateBackupCodesFailed")
	}

	// 开启事务
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	// 删除现有的TOTP记录（如果存在）
	if existingTotp != nil {
		err = tx.UserTotp.DeleteOneID(existingTotp.ID).Exec(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
	}

	// 创建新的TOTP记录
	_, err = tx.UserTotp.Create().
		SetUserID(userID).
		SetSecretKey(key.Secret()).
		SetBackupCodes(string(backupCodesJSON)).
		SetState(false). // 初始状态为未启用，需要验证后才启用
		SetIsVerified(false).
		SetIssuer(in.Issuer).
		Save(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 返回设置响应
	return &core.TotpSetupResponse{
		SecretKey:       key.Secret(),
		QrCodeContent:   key.URL(), // 返回TOTP URI
		BackupCodes:     backupCodes,
		DeviceName:      fmt.Sprintf("%s@%s", userInfo.Username, in.Domain),
		Issuer:          in.Issuer,
	}, nil
}

// generateBackupCodes 生成备用恢复码
func (l *EnableTotpLogic) generateBackupCodes(count int) []string {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		codes[i] = l.generateRandomCode(8)
	}
	return codes
}

// generateRandomCode 生成随机码
func (l *EnableTotpLogic) generateRandomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
