package userservicelogic

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent/user"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/userutils"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户
func (l *UpdateUserLogic) UpdateUser(in *core.UserInfo) (*core.UserInfo, error) {
	// 验证必填字段
	if err := l.validateUpdate(in); err != nil {
		return nil, err
	}

	// 解析用户ID
	userID, err := uuid.Parse(*in.Id)
	if err != nil {
		return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 开启事务
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	// 构建更新查询
	updateQuery := tx.User.UpdateOneID(userID)

	// 设置可更新的字段
	if in.Username != nil {
		updateQuery.SetUsername(*in.Username)
	}
	if in.Email != nil {
		updateQuery.SetEmail(*in.Email)
	}
	if in.FullName != nil {
		updateQuery.SetFullName(*in.FullName)
	}
	if in.Mobile != nil {
		updateQuery.SetMobile(*in.Mobile)
	}
	if in.Avatar != nil {
		updateQuery.SetAvatar(*in.Avatar)
	}
	if in.UserDescription != nil {
		updateQuery.SetUserDescription(*in.UserDescription)
	}
	if in.State != nil {
		updateQuery.SetState(*in.State)
	}
	if in.DepartmentId != nil {
		updateQuery.SetDepartmentID(*in.DepartmentId)
	}

	// 处理密码更新（如果提供了密码）
	if in.PasswordHash != nil && *in.PasswordHash != "" {
		updateQuery.SetPasswordHash(l.hashPassword(*in.PasswordHash))
	}

	// 执行基本信息更新
	_, err = updateQuery.Save(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 处理角色关联更新
	if in.RoleIds != nil {
		// 先清除现有角色关联
		err = tx.User.UpdateOneID(userID).ClearRoles().Exec(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
		// 添加新的角色关联
		if len(in.RoleIds) > 0 {
			roleIDs := make([]uint32, len(in.RoleIds))
			copy(roleIDs, in.RoleIds)
			err = tx.User.UpdateOneID(userID).AddRoleIDs(roleIDs...).Exec(l.ctx)
			if err != nil {
				return nil, errorhandler.DBEntError(l.Logger, err, in)
			}
		}
	}

	// 处理职位关联更新
	if in.PositionIds != nil {
		// 先清除现有职位关联
		err = tx.User.UpdateOneID(userID).ClearPositions().Exec(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
		// 添加新的职位关联
		if len(in.PositionIds) > 0 {
			positionIDs := make([]uint32, len(in.PositionIds))
			copy(positionIDs, in.PositionIds)
			err = tx.User.UpdateOneID(userID).AddPositionIDs(positionIDs...).Exec(l.ctx)
			if err != nil {
				return nil, errorhandler.DBEntError(l.Logger, err, in)
			}
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 重新查询用户信息（包含关联数据）
	userWithEdges, err := l.svcCtx.DBEnt.User.Query().
		Where(user.IDEQ(userID)).
		WithRoles().
		WithPositions().
		WithTotp().
		Only(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return userutils.ConvertUserToUserInfo(userWithEdges), nil
}

// 验证更新用户的必填字段
func (l *UpdateUserLogic) validateUpdate(in *core.UserInfo) error {
	if in.Id == nil || *in.Id == "" {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 密码哈希处理
func (l *UpdateUserLogic) hashPassword(password string) string {
	// 生成盐值
	salt := make([]byte, 16)
	rand.Read(salt)

	// 使用SHA256进行哈希
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hasher.Write(salt)

	// 返回十六进制编码的哈希值
	return hex.EncodeToString(hasher.Sum(nil)) + ":" + hex.EncodeToString(salt)
}

