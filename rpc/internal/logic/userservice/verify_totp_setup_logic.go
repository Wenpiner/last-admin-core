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

type VerifyTotpSetupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyTotpSetupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyTotpSetupLogic {
	return &VerifyTotpSetupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 验证并确认TOTP设置
func (l *VerifyTotpSetupLogic) VerifyTotpSetup(in *core.VerifyTotpSetupRequest) (*core.TotpSetupConfirmResponse, error) {
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
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 检查TOTP是否已经启用
	if totpRecord.IsEnabled {
		return &core.TotpSetupConfirmResponse{
			Success: false,
			Message: "totp.alreadyEnabled",
		}, nil
	}

	// 验证TOTP代码
	valid := totp.Validate(in.TotpCode, totpRecord.SecretKey)
	if !valid {
		return &core.TotpSetupConfirmResponse{
			Success: false,
			Message: "totp.verifyFailed",
		}, nil
	}

	// 开启事务
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	// 更新TOTP记录，启用并验证
	updateQuery := tx.UserTotp.UpdateOneID(totpRecord.ID).
		SetIsEnabled(true).
		SetIsVerified(true).
		SetLastUsedAt(time.Now()).
		SetLastUsedCode(in.TotpCode)

	// 如果提供了设备名称，更新设备名称
	if in.DeviceName != nil && *in.DeviceName != "" {
		updateQuery.SetDeviceName(*in.DeviceName)
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

	return &core.TotpSetupConfirmResponse{
		Success:     true,
		Message:     "totp.setupSuccess",
		BackupCodes: []string{}, // 备用码在EnableTotp时已经返回
	}, nil
}
