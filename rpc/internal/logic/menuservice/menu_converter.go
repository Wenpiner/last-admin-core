package menuservicelogic

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/types/core"
)

// ConvertMenuToMenuInfo 将 Menu 实体转换为 MenuInfo
func ConvertMenuToMenuInfo(menu *ent.Menu) *core.MenuInfo {
	return &core.MenuInfo{
		Id:          &menu.ID,
		CreatedAt:   pointer.ToInt64Ptr(menu.CreatedAt.UnixMilli()),
		UpdatedAt:   pointer.ToInt64Ptr(menu.UpdatedAt.UnixMilli()),
		MenuCode:    &menu.MenuCode,
		MenuName:    &menu.MenuName,
		ParentId:    &menu.ParentID,
		MenuPath:    menu.MenuPath,
		State:       &menu.State,
		Sort:        &menu.Sort,
		MenuType:    &menu.MenuType,
		Description: menu.Description,
		Component:   menu.Component,
		Redirect:    menu.Redirect,
		ServiceName: menu.ServiceName,
		MenuLevel:   pointer.ToUint32Ptr(uint32(menu.MenuLevel)),
		Permission:  menu.Permission,
		Meta: &core.MenuMeta{
			Title:        &menu.MenuName, // menu_name 对应 meta.title
			Icon:         menu.Icon,
			IsHidden:     menu.IsHidden,
			IsBreadcrumb: menu.IsBreadcrumb,
			IsCache:      menu.IsCache,
			IsTab:        menu.IsTab,
			IsAffix:      menu.IsAffix,
			FrameSrc:     menu.FrameSrc,
			Link:         menu.Link,
		},
	}
}
