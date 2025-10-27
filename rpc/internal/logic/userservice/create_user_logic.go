package userservicelogic

import (
	"context"

	"github.com/google/uuid"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-common/utils/encrypt"
	"github.com/wenpiner/last-admin-core/rpc/ent/role"
	"github.com/wenpiner/last-admin-core/rpc/ent/user"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/userutils"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建用户
func (l *CreateUserLogic) CreateUser(in *core.UserInfo) (*core.UserInfo, error) {
	// 验证必填字段
	if err := l.validateCreate(in); err != nil {
		return nil, err
	}

	// 开启事务
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	// 生成用户ID
	userID := uuid.New()

	// 创建用户基本信息
	createQuery := tx.User.Create().
		SetID(userID).
		SetUsername(*in.Username).
		SetPasswordHash(l.hashPassword(*in.PasswordHash)).
		SetState(l.getStateValue(in.State))

	// 设置可选字段
	if in.Email != nil {
		createQuery.SetEmail(*in.Email)
	}
	if in.FullName != nil {
		createQuery.SetFullName(*in.FullName)
	}
	if in.Mobile != nil {
		createQuery.SetMobile(*in.Mobile)
	}
	if in.Avatar != nil {
		createQuery.SetAvatar(*in.Avatar)
	}
	if in.UserDescription != nil {
		createQuery.SetUserDescription(*in.UserDescription)
	}

	if in.DepartmentId != nil {
		createQuery.SetDepartmentID(*in.DepartmentId)
	}

	// 创建用户
	_, err = createQuery.Save(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	var roleIds []uint32
	// 处理角色关联
	if len(in.RoleIds) > 0 {
		roleIds = append(roleIds, in.RoleIds...)
	}

	if len(in.RoleValues) > 0 {
		l.Infow("用户注册", logx.Field("roleValues", in.RoleValues))
		// 查询角色ID
		roles, err := tx.Role.Query().Where(role.RoleCodeIn(in.RoleValues...)).All(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
		for _, role := range roles {
			roleIds = append(roleIds, role.ID)
		}
	}

	if len(roleIds) > 0 {
		err = tx.User.UpdateOneID(userID).AddRoleIDs(roleIds...).Exec(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
	}

	// 处理职位关联
	if len(in.PositionIds) > 0 {
		positionIDs := make([]uint32, len(in.PositionIds))
		copy(positionIDs, in.PositionIds)
		err = tx.User.UpdateOneID(userID).AddPositionIDs(positionIDs...).Exec(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
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

// 验证创建用户的必填字段
func (l *CreateUserLogic) validateCreate(in *core.UserInfo) error {
	if in.Username == nil || *in.Username == "" {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	if in.PasswordHash == nil || *in.PasswordHash == "" {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	if in.DepartmentId == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 密码哈希处理
func (l *CreateUserLogic) hashPassword(password string) string {
	return encrypt.BcryptEncrypt(password)
}

// 获取状态值，默认为 true (启用)
func (l *CreateUserLogic) getStateValue(state *bool) bool {
	if state != nil {
		return *state
	}
	return true
}
