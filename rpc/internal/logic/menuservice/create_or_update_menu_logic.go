package menuservicelogic

import (
	"context"

	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdateMenuLogic {
	return &CreateOrUpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建或更新菜单
func (l *CreateOrUpdateMenuLogic) CreateOrUpdateMenu(in *core.MenuInfo) (*core.MenuInfo, error) {
	// 开启事务，并先进行检查是否存在，如果存在则进行更新否则进行创建
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	var menu *ent.Menu
	if in.Id == nil {
		// 新增,验证必填参数可用性
		if err := l.validateCreate(in); err != nil {
			return nil, err
		}

		// 计算菜单层级
		menuLevel, err := l.calculateMenuLevel(tx, in.ParentId)
		if err != nil {
			return nil, err
		}

		menu, err = tx.Menu.Create().
			SetMenuCode(pointer.GetString(in.MenuCode)).
			SetMenuName(pointer.GetString(in.MenuName)).
			SetNillableParentID(in.ParentId).
			SetMenuLevel(menuLevel).
			SetNillableMenuPath(in.MenuPath).
			SetNillableComponent(in.Component).
			SetNillableRedirect(in.Redirect).
			SetNillableServiceName(in.ServiceName).
			SetMenuType(pointer.GetString(in.MenuType)).
			SetNillableFrameSrc(l.getFrameSrcFromMeta(in.Meta)).
			SetNillableDescription(in.Description).
			SetState(pointer.GetBool(in.State)).
			SetSort(l.getSortValue(in.Sort)).
			SetNillableIcon(l.getIconFromMeta(in.Meta)).
			SetNillablePermission(in.Permission).
			SetIsHidden(l.getIsHiddenFromMeta(in.Meta)).
			SetIsBreadcrumb(l.getIsBreadcrumbFromMeta(in.Meta)).
			SetIsCache(l.getIsCacheFromMeta(in.Meta)).
			SetIsTab(l.getIsTabFromMeta(in.Meta)).
			SetIsAffix(l.getIsAffixFromMeta(in.Meta)).
			Save(l.ctx)
	} else {
		// 更新
		if err := l.validateUpdate(in); err != nil {
			return nil, err
		}

		updateQuery := tx.Menu.UpdateOneID(pointer.GetUint32(in.Id)).
			SetNillableMenuCode(in.MenuCode).
			SetNillableMenuName(in.MenuName).
			SetNillableMenuPath(in.MenuPath).
			SetNillableComponent(in.Component).
			SetNillableRedirect(in.Redirect).
			SetNillableServiceName(in.ServiceName).
			SetNillableMenuType(in.MenuType).
			SetNillableFrameSrc(l.getFrameSrcFromMeta(in.Meta)).
			SetNillableDescription(in.Description).
			SetNillableState(in.State).
			SetNillableSort(in.Sort).
			SetNillableIcon(l.getIconFromMeta(in.Meta)).
			SetNillablePermission(in.Permission).
			SetNillableIsHidden(l.getNillableIsHiddenFromMeta(in.Meta)).
			SetNillableIsBreadcrumb(l.getNillableIsBreadcrumbFromMeta(in.Meta)).
			SetNillableIsCache(l.getNillableIsCacheFromMeta(in.Meta)).
			SetNillableIsTab(l.getNillableIsTabFromMeta(in.Meta)).
			SetNillableIsAffix(l.getNillableIsAffixFromMeta(in.Meta))

		// 如果更新了父菜单ID，需要重新计算层级
		if in.ParentId != nil {
			menuLevel, err := l.calculateMenuLevel(tx, in.ParentId)
			if err != nil {
				return nil, err
			}
			updateQuery = updateQuery.SetMenuLevel(menuLevel)
		}

		menu, err = updateQuery.Save(l.ctx)
	}
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return ConvertMenuToMenuInfo(menu), nil
}

// 验证新增参数可用性
func (l *CreateOrUpdateMenuLogic) validateCreate(in *core.MenuInfo) error {
	if in.MenuCode == nil || in.MenuName == nil || in.MenuType == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 所属服务为必填项，不管是哪种类型
	if in.ServiceName == nil || *in.ServiceName == "" {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	
	// 根据菜单类型进行特定验证
	menuType := pointer.GetString(in.MenuType)
	switch menuType {
	case "directory":
		// 目录时应该保证路由地址、图标
		if in.MenuPath == nil || *in.MenuPath == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.Meta == nil || in.Meta.Icon == nil || *in.Meta.Icon == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
	case "menu":
		// 菜单时应该保证路由地址、组件地址
		if in.MenuPath == nil || *in.MenuPath == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.Component == nil || *in.Component == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
		if in.Meta == nil || in.Meta.Icon == nil || *in.Meta.Icon == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
	case "button":
		// 按钮/页面元素时保证权限标识
		if in.Permission == nil || *in.Permission == "" {
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
	default:
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	return nil
}

// 验证更新参数可用性
func (l *CreateOrUpdateMenuLogic) validateUpdate(in *core.MenuInfo) error {
	if in.Id == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 如果更新了所属服务，不能为空
	if in.ServiceName != nil && *in.ServiceName == "" {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 保证名称
	if in.MenuName == nil || *in.MenuName == "" {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	if in.MenuCode == nil || *in.MenuCode == "" {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}

	// 如果更新了菜单类型，需要进行相应验证
	if in.MenuType != nil {
		menuType := pointer.GetString(in.MenuType)
		switch menuType {
		case "directory":
			// 目录时应该保证路由地址、图标
			if in.MenuPath != nil && *in.MenuPath == "" {
				return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
			}
			if in.Meta != nil && in.Meta.Icon != nil && *in.Meta.Icon == "" {
				return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
			}
		case "menu":
			// 菜单时应该保证路由地址、组件地址、图标
			if in.MenuPath != nil && *in.MenuPath == "" {
				return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
			}
			if in.Component != nil && *in.Component == "" {
				return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
			}
			if in.Meta != nil && in.Meta.Icon != nil && *in.Meta.Icon == "" {
				return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
			}
		case "button":
			// 按钮/页面元素时保证权限标识
			if in.Permission != nil && *in.Permission == "" {
				return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
			}
		default:
			return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
	}

	return nil
}

// 计算菜单层级
func (l *CreateOrUpdateMenuLogic) calculateMenuLevel(tx *ent.Tx, parentId *uint32) (uint16, error) {
	// 如果没有父菜单ID或父菜单ID为0，则为顶级菜单，层级为0
	if parentId == nil || *parentId == 0 {
		return 0, nil
	}

	// 查询父菜单信息
	parentMenu, err := tx.Menu.Get(l.ctx, *parentId)
	if err != nil {
		return 0, errorhandler.DBEntError(l.Logger, err, parentId)
	}

	// 父菜单层级 + 1
	return parentMenu.MenuLevel + 1, nil
}

// 获取排序值
func (l *CreateOrUpdateMenuLogic) getSortValue(sort *int32) int32 {
	if sort == nil {
		return 0
	}
	return *sort
}

// 从 Meta 中获取图标
func (l *CreateOrUpdateMenuLogic) getIconFromMeta(meta *core.MenuMeta) *string {
	if meta == nil {
		return nil
	}
	return meta.Icon
}


// 从 Meta 中获取是否隐藏
func (l *CreateOrUpdateMenuLogic) getIsHiddenFromMeta(meta *core.MenuMeta) bool {
	if meta == nil {
		return false
	}
	return pointer.GetBool(meta.IsHidden)
}

// 从 Meta 中获取是否显示在面包屑
func (l *CreateOrUpdateMenuLogic) getIsBreadcrumbFromMeta(meta *core.MenuMeta) bool {
	if meta == nil {
		return true // 默认显示在面包屑
	}
	return pointer.GetBool(meta.IsBreadcrumb)
}

// 从 Meta 中获取是否缓存
func (l *CreateOrUpdateMenuLogic) getIsCacheFromMeta(meta *core.MenuMeta) bool {
	if meta == nil {
		return false
	}
	return pointer.GetBool(meta.IsCache)
}

// 从 Meta 中获取是否显示在标签栏
func (l *CreateOrUpdateMenuLogic) getIsTabFromMeta(meta *core.MenuMeta) bool {
	if meta == nil {
		return false
	}
	return pointer.GetBool(meta.IsTab)
}

// 从 Meta 中获取是否固定在标签栏
func (l *CreateOrUpdateMenuLogic) getIsAffixFromMeta(meta *core.MenuMeta) bool {
	if meta == nil {
		return false
	}
	return pointer.GetBool(meta.IsAffix)
}

// 以下是用于更新操作的 Nillable 版本函数

// 从 Meta 中获取是否隐藏（可空）
func (l *CreateOrUpdateMenuLogic) getNillableIsHiddenFromMeta(meta *core.MenuMeta) *bool {
	if meta == nil || meta.IsHidden == nil {
		return nil
	}
	return meta.IsHidden
}

// 从 Meta 中获取是否显示在面包屑（可空）
func (l *CreateOrUpdateMenuLogic) getNillableIsBreadcrumbFromMeta(meta *core.MenuMeta) *bool {
	if meta == nil || meta.IsBreadcrumb == nil {
		return nil
	}
	return meta.IsBreadcrumb
}

// 从 Meta 中获取是否缓存（可空）
func (l *CreateOrUpdateMenuLogic) getNillableIsCacheFromMeta(meta *core.MenuMeta) *bool {
	if meta == nil || meta.IsCache == nil {
		return nil
	}
	return meta.IsCache
}

// 从 Meta 中获取是否显示在标签栏（可空）
func (l *CreateOrUpdateMenuLogic) getNillableIsTabFromMeta(meta *core.MenuMeta) *bool {
	if meta == nil || meta.IsTab == nil {
		return nil
	}
	return meta.IsTab
}

// 从 Meta 中获取是否固定在标签栏（可空）
func (l *CreateOrUpdateMenuLogic) getNillableIsAffixFromMeta(meta *core.MenuMeta) *bool {
	if meta == nil || meta.IsAffix == nil {
		return nil
	}
	return meta.IsAffix
}

// 从 Meta 中获取内嵌 iframe 的 URL
func (l *CreateOrUpdateMenuLogic) getFrameSrcFromMeta(meta *core.MenuMeta) *string {
	if meta == nil || meta.FrameSrc == nil {
		return nil
	}
	return meta.FrameSrc
}


