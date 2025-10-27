package userservicelogic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent/usertotp"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type UseBackupCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseBackupCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseBackupCodeLogic {
	return &UseBackupCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 使用备用恢复码
func (l *UseBackupCodeLogic) UseBackupCode(in *core.UseBackupCodeRequest) (*core.BaseResponse, error) {
	// 解析用户ID
	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 查询用户的TOTP记录
	totpRecord, err := l.svcCtx.DBEnt.UserTotp.Query().
		Where(usertotp.UserIDEQ(userID)).
		First(l.ctx)
	if err != nil {
		return &core.BaseResponse{
			Message: "totp.notEnabled",
		}, nil
	}

	// 检查TOTP是否已启用和验证
	if !totpRecord.IsEnabled || !totpRecord.IsVerified {
		return &core.BaseResponse{
			Message: "totp.notEnabled",
		}, nil
	}

	// 检查账户是否被锁定
	if totpRecord.LockedUntil != nil && time.Now().Before(*totpRecord.LockedUntil) {
		return &core.BaseResponse{
			Message: "totp.accountLocked",
		}, nil
	}

	// 解析备用恢复码
	if totpRecord.BackupCodes == nil {
		return &core.BaseResponse{
			Message: "totp.noBackupCodes",
		}, nil
	}

	var backupCodes []string
	err = json.Unmarshal([]byte(*totpRecord.BackupCodes), &backupCodes)
	if err != nil {
		return &core.BaseResponse{
			Message: "totp.backupCodesError",
		}, nil
	}

	// 查找并验证备用码
	var foundIndex = -1
	for i, code := range backupCodes {
		if code == in.BackupCode {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return &core.BaseResponse{
			Message: "totp.invalidBackupCode",
		}, nil
	}

	// 开启事务
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	// 移除已使用的备用码
	backupCodes = append(backupCodes[:foundIndex], backupCodes[foundIndex+1:]...)
	newBackupCodesJSON, err := json.Marshal(backupCodes)
	if err != nil {
		return nil, errorx.NewInternalError("totp.updateRecordFailed")
	}

	// 更新TOTP记录
	_, err = tx.UserTotp.UpdateOneID(totpRecord.ID).
		SetBackupCodes(string(newBackupCodesJSON)).
		SetFailureCount(0).
		SetLastUsedAt(time.Now()).
		ClearLockedUntil().
		Save(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return &core.BaseResponse{
		Message: "totp.backupCodeUsed",
	}, nil
}
