package operationlogservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOperationLogLogic {
	return &CreateOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOperationLogLogic) CreateOperationLog(in *core.OperationLogInfo) (*core.BaseResponse, error) {
	// todo: add your logic here and delete this line

	return &core.BaseResponse{}, nil
}
