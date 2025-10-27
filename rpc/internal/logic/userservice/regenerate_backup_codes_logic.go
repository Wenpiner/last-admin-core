package userservicelogic

import (
	"context"
	"crypto/rand"
	"encoding/json"

	"github.com/google/uuid"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent/usertotp"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegenerateBackupCodesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegenerateBackupCodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegenerateBackupCodesLogic {
	return &RegenerateBackupCodesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 重新生成备用恢复码
func (l *RegenerateBackupCodesLogic) RegenerateBackupCodes(in *core.UUIDRequest) (*core.BackupCodesResponse, error) {
	// 解析用户ID
	userID, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 查询用户的TOTP记录
	totpRecord, err := l.svcCtx.DBEnt.UserTotp.Query().
		Where(usertotp.UserIDEQ(userID)).
		First(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 检查TOTP是否已启用
	if !totpRecord.IsEnabled {
		return &core.BackupCodesResponse{
			BackupCodes: []string{},
			Message:     "totp.notEnabled",
		}, nil
	}

	// 生成新的备用恢复码
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

	// 更新备用恢复码
	_, err = tx.UserTotp.UpdateOneID(totpRecord.ID).
		SetBackupCodes(string(backupCodesJSON)).
		Save(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return &core.BackupCodesResponse{
		BackupCodes: backupCodes,
		Message:     "totp.backupCodesRegenerated",
	}, nil
}

// generateBackupCodes 生成备用恢复码
func (l *RegenerateBackupCodesLogic) generateBackupCodes(count int) []string {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		codes[i] = l.generateRandomCode(8)
	}
	return codes
}

// generateRandomCode 生成随机码
func (l *RegenerateBackupCodesLogic) generateRandomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
