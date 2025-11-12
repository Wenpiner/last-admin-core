package userservicelogic

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent/usertotp"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTotpStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTotpStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTotpStatusLogic {
	return &GetTotpStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户TOTP状态
func (l *GetTotpStatusLogic) GetTotpStatus(in *core.UUIDRequest) (*core.TotpStatusResponse, error) {
	// 解析用户ID
	userID, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 查询用户的TOTP记录
	totpRecord, err := l.svcCtx.DBEnt.UserTotp.Query().
		Where(usertotp.IDEQ(userID)).
		First(l.ctx)
	if err != nil {
		// 如果没有找到TOTP记录，返回未启用状态
		return &core.TotpStatusResponse{
			State:        false,
			IsVerified:       false,
			BackupCodesCount: 0,
		}, nil
	}

	// 计算备用码数量
	backupCodesCount := int32(0)
	if totpRecord.BackupCodes != nil {
		var backupCodes []string
		if err := json.Unmarshal([]byte(*totpRecord.BackupCodes), &backupCodes); err == nil {
			backupCodesCount = int32(len(backupCodes))
		}
	}

	// 构建响应
	response := &core.TotpStatusResponse{
		State:        totpRecord.State,
		IsVerified:       totpRecord.IsVerified,
		BackupCodesCount: backupCodesCount,
	}

	// 设置可选字段
	if totpRecord.DeviceName != nil {
		response.DeviceName = totpRecord.DeviceName
	}

	if totpRecord.LastUsedAt != nil {
		timestamp := totpRecord.LastUsedAt.UnixMilli()
		response.LastUsedAt = &timestamp
	}
	return response, nil
}
