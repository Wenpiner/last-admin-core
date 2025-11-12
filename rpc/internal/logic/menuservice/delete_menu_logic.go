package menuservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/menu"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除菜单
func (l *DeleteMenuLogic) DeleteMenu(in *core.ID32Request) (*core.BaseResponse, error) {
	// 判断是否存在子菜单
	hasChildren, err := l.svcCtx.DBEnt.Menu.Query().Where(menu.ParentIDEQ(in.Id)).Exist(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	if hasChildren {
		return nil, errorx.NewInvalidArgumentError("menu.hasChildren")
	}

	err = l.svcCtx.DBEnt.Menu.DeleteOneID(in.Id).Exec(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return &core.BaseResponse{
		Message: "common.deleteSuccess",
	}, nil
}
