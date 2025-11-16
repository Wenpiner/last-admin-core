package operationlogservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperationLogLogic {
	return &ListOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListOperationLogLogic) ListOperationLog(in *core.OperationLogListRequest) (*core.OperationLogListResponse, error) {
	// todo: add your logic here and delete this line

	return &core.OperationLogListResponse{}, nil
}
