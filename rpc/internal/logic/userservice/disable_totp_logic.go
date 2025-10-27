package userservicelogic

import (
	"context"

	"github.com/google/uuid"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent/usertotp"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type DisableTotpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDisableTotpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisableTotpLogic {
	return &DisableTotpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 禁用TOTP
func (l *DisableTotpLogic) DisableTotp(in *core.DisableTotpRequest) (*core.BaseResponse, error) {
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

	// 开启事务
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	// 删除TOTP记录
	err = tx.UserTotp.DeleteOneID(totpRecord.ID).Exec(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return &core.BaseResponse{
		Message: "totp.disableSuccess",
	}, nil
}
