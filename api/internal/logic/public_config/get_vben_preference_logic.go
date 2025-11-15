package public_config

import (
	"context"
	"encoding/json"

	"github.com/wenpiner/last-admin-common/enums"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetVbenPreferenceLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取VBen Preference配置
func NewGetVbenPreferenceLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetVbenPreferenceLogic {
	return &GetVbenPreferenceLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetVbenPreferenceLogic) GetVbenPreference() (resp *types.VbenPreference, err error) {
	result, err := l.svcCtx.ConfigurationRpc.GetConfiguration(l.ctx,&core.StringRequest{
		Value: enums.ConfigurationVBenPreference,
	})
	if err != nil {
		return nil, errorx.NewApiError(errorx.CodeInternalError, err.Error())
	}

	// 将 result.Value 转 map[string]interface{}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(result.Value), &data)
	if err != nil {
		return nil, errorx.NewApiError(errorx.CodeInternalError, err.Error())
	}

	resp = &types.VbenPreference{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: data,
	}

	return
}
