package userservicelogic

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent/usertotp"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyTotpCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyTotpCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyTotpCodeLogic {
	return &VerifyTotpCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 验证TOTP代码（用于登录）
func (l *VerifyTotpCodeLogic) VerifyTotpCode(in *core.VerifyTotpCodeRequest) (*core.VerifyTotpCodeResponse, error) {
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
		return &core.VerifyTotpCodeResponse{
			IsValid: false,
			Message: "totp.notEnabled",
		}, nil
	}

	// 检查TOTP是否已启用和验证
	if !totpRecord.IsEnabled || !totpRecord.IsVerified {
		return &core.VerifyTotpCodeResponse{
			IsValid: false,
			Message: "totp.notEnabled",
		}, nil
	}

	// 检查账户是否被锁定
	if totpRecord.LockedUntil != nil && time.Now().Before(*totpRecord.LockedUntil) {
		lockedUntilTimestamp := totpRecord.LockedUntil.UnixMilli()
		return &core.VerifyTotpCodeResponse{
			IsValid:     false,
			Message:     "totp.accountLocked",
			LockedUntil: &lockedUntilTimestamp,
		}, nil
	}

	// 防重放攻击：检查是否使用了相同的验证码
	if totpRecord.LastUsedCode != nil && *totpRecord.LastUsedCode == in.TotpCode {
		return &core.VerifyTotpCodeResponse{
			IsValid: false,
			Message: "totp.codeAlreadyUsed",
		}, nil
	}

	// 验证TOTP代码
	valid := totp.Validate(in.TotpCode, totpRecord.SecretKey)

	// 开启事务
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	if valid {
		// 验证成功，重置失败次数并更新使用记录
		_, err = tx.UserTotp.UpdateOneID(totpRecord.ID).
			SetFailureCount(0).
			SetLastUsedAt(time.Now()).
			SetLastUsedCode(in.TotpCode).
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

		return &core.VerifyTotpCodeResponse{
			IsValid: true,
			Message: "totp.verifySuccess",
		}, nil
	} else {
		// 验证失败，增加失败次数
		newFailureCount := totpRecord.FailureCount + 1
		updateQuery := tx.UserTotp.UpdateOneID(totpRecord.ID).
			SetFailureCount(newFailureCount)

		// 如果失败次数达到5次，锁定账户30分钟
		const maxFailures = 5
		const lockDuration = 30 * time.Minute
		var lockedUntil *int64

		if newFailureCount >= maxFailures {
			lockTime := time.Now().Add(lockDuration)
			updateQuery.SetLockedUntil(lockTime)
			lockedUntil = &[]int64{lockTime.UnixMilli()}[0]
		}

		_, err = updateQuery.Save(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}

		// 提交事务
		err = tx.Commit()
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}

		// 计算剩余尝试次数
		remainingAttempts := int32(maxFailures - newFailureCount)
		if remainingAttempts < 0 {
			remainingAttempts = 0
		}

		return &core.VerifyTotpCodeResponse{
			IsValid:           false,
			Message:           "totp.verifyFailed",
			RemainingAttempts: &remainingAttempts,
			LockedUntil:       lockedUntil,
		}, nil
	}
}
