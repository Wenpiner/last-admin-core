package initservicelogic

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/bsm/redislock"
	"github.com/wenpiner/last-admin-common/enums"
	"github.com/wenpiner/last-admin-common/utils/encrypt"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type InitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitLogic {
	return &InitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitLogic) Init(in *core.EmptyRequest) (*core.BaseResponse, error) {
	// 由于创建时间不确定，所以使用自定义ctx
	ctx := context.Background()
	// Redis 获取创建锁
	locker := redislock.New(l.svcCtx.Redis)
	lockKey := "last-admin-core:init"
	lock, err := locker.Obtain(ctx, lockKey, 10*time.Minute, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		// 正在初始化中
		return nil, errorx.NewInternalError("init.pending")
	} else if err != nil {
		return nil, errorx.NewInternalError("init.failed")
	}
	defer lock.Release(ctx)

	// initialize table structure
	if err := l.svcCtx.DBEnt.Schema.Create(l.ctx, schema.WithForeignKeys(false), schema.WithDropColumn(true),
		schema.WithDropIndex(true)); err != nil {
		logx.Errorw("初始化数据库失败", logx.Field("detail", err.Error()))
		_ = l.svcCtx.Redis.Set(l.ctx, "INIT:DATABASE:ERROR", err.Error(), 300*time.Second).Err()
		return nil, errorx.NewInternalError(err.Error())
	}

	// 初始化角色信息
	if err := l.initRole(); err != nil {
		return nil, err
	}

	// 初始化用户信息
	if err := l.initUser(); err != nil {
		return nil, err
	}

	// 初始化菜单
	if err := l.initMenu(); err != nil {
		return nil, err
	}

	// 初始化菜单和角色关联
	if err := l.initRoleMenu(); err != nil {
		return nil, err
	}

	if err := l.initBaseApi(); err != nil {
		return nil, err
	}

	if err := l.initCasbin(); err != nil {
		return nil, err
	}

	if err := l.initOauthProvider(); err != nil {
		return nil, err
	}

	if err := l.initDepartment(); err != nil {
		return nil, err
	}

	if err := l.initPosition(); err != nil {
		return nil, err
	}

	return &core.BaseResponse{}, nil
}

// 初始化casbin配置
func (l *InitLogic) initCasbin() error {

	// 查询出来所有API信息
	apis, err := l.svcCtx.DBEnt.API.Query().All(l.ctx)
	if err != nil {
		logx.Errorw("查询API信息失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}

	var policies [][]string
	for _, api := range apis {
		policies = append(policies, []string{enums.DefaultRoleValue, api.Path, api.Method})
	}

	// // 查询所有菜单
	// menus, err := l.svcCtx.DBEnt.Menu.Query().All(l.ctx)
	// if err != nil {
	// 	logx.Errorw("查询菜单信息失败", logx.Field("detail", err.Error()))
	// 	return errorx.NewInternalError(err.Error())
	// }
	// for _, menu := range menus {
	// 	policies = append(policies, []string{enums.DefaultRoleValue, pointer.GetString(menu.MenuPath), menu.MenuType})
	// }

	var oldPolicies [][]string
	oldPolicies, err = l.svcCtx.Casbin.GetFilteredPolicy(0, enums.DefaultRoleValue)
	if err != nil {
		logx.Errorw("查询 Casbin 策略失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}

	if len(oldPolicies) != 0 {
		removeResult, err := l.svcCtx.Casbin.RemoveFilteredPolicy(0, enums.DefaultRoleValue)
		if err != nil {
			logx.Errorw("删除 Casbin 策略失败", logx.Field("detail", err.Error()))
			return errorx.NewInternalError(err.Error())
		}
		if !removeResult {
			logx.Errorw("删除 Casbin 策略失败", logx.Field("detail", "删除失败"))
			return errorx.NewInternalError("删除 Casbin 策略失败")
		}
	}

	// 添加新的策略
	if result, err := l.svcCtx.Casbin.AddPolicies(policies); err != nil || !result {
		logx.Errorw("添加 Casbin 策略失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}

	return nil
}

// 初始化基础API配置
func (l *InitLogic) initBaseApi() error {
	var apis []*ent.APICreate

	// Auth
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Auth").
		SetMethod("POST").
		SetPath("/auth/login").
		SetServiceName("core").
		SetName("账号密码登录").SetIsRequired(true))
	// /auth/oauth/callback
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Auth").
		SetMethod("GET").
		SetPath("/auth/oauth/callback").
		SetServiceName("core").
		SetName("Oauth 回调").SetIsRequired(true))
	// /auth/oauth/callback
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Auth").
		SetMethod("POST").
		SetPath("/auth/oauth/login").
		SetServiceName("core").
		SetName("Oauth 登录").SetIsRequired(true))
	// /auth/oauth/callback
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Auth").
		SetMethod("POST").
		SetPath("/auth/register").
		SetServiceName("core").
		SetName("注册用户").SetIsRequired(true))
	// /auth/codes
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Auth").
		SetMethod("GET").
		SetPath("/auth/codes").
		SetServiceName("core").
		SetName("获取用户权限码(菜单按钮级别)").SetIsRequired(true))

	// Captcha
	// /captcha/email
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Captcha").
		SetMethod("POST").
		SetPath("/captcha/email").
		SetServiceName("core").
		SetName("发送邮箱验证码").SetIsRequired(true))

	// /captcha/generate
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Captcha").
		SetMethod("GET").
		SetPath("/captcha/generate").
		SetServiceName("core").
		SetName("生成验证码").SetIsRequired(true))

	// /captcha/sms
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Captcha").
		SetMethod("POST").
		SetPath("/captcha/sms").
		SetServiceName("core").
		SetName("发送短信验证码").SetIsRequired(true))

	// User
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("User").
		SetMethod("GET").
		SetPath("/user/info").
		SetServiceName("core").
		SetName("获取用户信息").SetIsRequired(true))

	// Menu
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Menu").
		SetMethod("GET").
		SetPath("/menu/all").
		SetServiceName("core").
		SetName("获取用户角色当前所有菜单").SetIsRequired(true))

	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Menu").
		SetMethod("GET").
		SetPath("/menu/all-menus").
		SetServiceName("core").
		SetName("获取所有菜单").SetIsRequired(true))

	// delete
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Menu").
		SetMethod("POST").
		SetPath("/menu/delete").
		SetServiceName("core").
		SetName("删除菜单").SetIsRequired(false))

	// update
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Menu").
		SetMethod("PUT").
		SetPath("/menu/update").
		SetServiceName("core").
		SetName("更新菜单").SetIsRequired(false))

	// Api
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Api").
		SetMethod("POST").
		SetPath("/api/list").
		SetServiceName("core").
		SetName("获取API列表").SetIsRequired(false))

	// /api/all
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Api").
		SetMethod("GET").
		SetPath("/api/all").
		SetServiceName("core").
		SetName("获取所有API").SetIsRequired(false))

	// create or update
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Api").
		SetMethod("POST").
		SetPath("/api/createOrUpdate").
		SetServiceName("core").
		SetName("创建或更新API").SetIsRequired(false))

	// delete
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Api").
		SetMethod("POST").
		SetPath("/api/delete").
		SetServiceName("core").
		SetName("删除API").SetIsRequired(false))

	// Department
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Department").
		SetMethod("POST").
		SetPath("/department/createOrUpdate").
		SetServiceName("core").
		SetName("创建或更新部门").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Department").
		SetMethod("POST").
		SetPath("/department/delete").
		SetServiceName("core").
		SetName("删除部门").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Department").
		SetMethod("POST").
		SetPath("/department/list").
		SetServiceName("core").
		SetName("获取部门列表").SetIsRequired(false))

	// user
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("User").
		SetMethod("POST").
		SetPath("/user/list").
		SetServiceName("core").
		SetName("获取用户列表").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("User").
		SetMethod("POST").
		SetPath("/user/createOrUpdate").
		SetServiceName("core").
		SetName("创建或更新用户").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("User").
		SetMethod("POST").
		SetPath("/user/delete").
		SetServiceName("core").
		SetName("删除用户").SetIsRequired(false))

	// role
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Role").
		SetMethod("POST").
		SetPath("/role/createOrUpdate").
		SetServiceName("core").
		SetName("创建或更新角色").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Role").
		SetMethod("POST").
		SetPath("/role/delete").
		SetServiceName("core").
		SetName("删除角色").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Role").
		SetMethod("POST").
		SetPath("/role/list").
		SetServiceName("core").
		SetName("获取角色列表").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Role").
		SetMethod("POST").
		SetPath("/role/assign/menu").
		SetServiceName("core").
		SetName("为角色分配菜单").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Role").
		SetMethod("POST").
		SetPath("/role/assign/api").
		SetServiceName("core").
		SetName("为角色分配API").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Role").
		SetMethod("POST").
		SetPath("/role/get/menu").
		SetServiceName("core").
		SetName("获取角色菜单").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Role").
		SetMethod("POST").
		SetPath("/role/get/api").
		SetServiceName("core").
		SetName("获取角色API").SetIsRequired(false))

	// Position
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Position").
		SetMethod("POST").
		SetPath("/position/createOrUpdate").
		SetServiceName("core").
		SetName("创建或更新岗位").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Position").
		SetMethod("POST").
		SetPath("/position/delete").
		SetServiceName("core").
		SetName("删除岗位").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Position").
		SetMethod("POST").
		SetPath("/position/list").
		SetServiceName("core").
		SetName("获取岗位列表").SetIsRequired(false))

	// Dict
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Dict").
		SetMethod("POST").
		SetPath("/dict/createOrUpdate").
		SetServiceName("core").
		SetName("创建或更新字典").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Dict").
		SetMethod("POST").
		SetPath("/dict/delete").
		SetServiceName("core").
		SetName("删除字典").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Dict").
		SetMethod("POST").
		SetPath("/dict/list").
		SetServiceName("core").
		SetName("获取字典列表").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Dict").
		SetMethod("POST").
		SetPath("/dict/get").
		SetServiceName("core").
		SetName("获取字典").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Dict").
		SetMethod("POST").
		SetPath("/dict/item/createOrUpdate").
		SetServiceName("core").
		SetName("创建或更新字典子项").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Dict").
		SetMethod("POST").
		SetPath("/dict/item/delete").
		SetServiceName("core").
		SetName("删除字典子项").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Dict").
		SetMethod("POST").
		SetPath("/dict/item/list").
		SetServiceName("core").
		SetName("获取字典子项列表").SetIsRequired(false))
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Dict").
		SetMethod("POST").
		SetPath("/dict/item/get").
		SetServiceName("core").
		SetName("获取字典子项").SetIsRequired(false))

	// Oauth
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Oauth").
		SetMethod("POST").
		SetPath("/oauth/list").
		SetServiceName("core").
		SetName("Oauth列表").SetIsRequired(false))
	//  /oauth/createOrUpdate
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Oauth").
		SetMethod("POST").
		SetPath("/oauth/createOrUpdate").
		SetServiceName("core").
		SetName("创建或更新Oauth").SetIsRequired(false))
	//  /oauth/delete
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Oauth").
		SetMethod("POST").
		SetPath("/oauth/delete").
		SetServiceName("core").
		SetName("删除Oauth").SetIsRequired(false))

	// Token
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Token").
		SetMethod("POST").
		SetPath("/token/list").
		SetServiceName("core").
		SetName("获取Token列表").SetIsRequired(false))

	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Token").
		SetMethod("POST").
		SetPath("/token/clean").
		SetServiceName("core").
		SetName("清理过期Token").SetIsRequired(false))
		
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Token").
		SetMethod("POST").
		SetPath("/token/block").
		SetServiceName("core").
		SetName("拉黑用户Token").SetIsRequired(false))
		
	apis = append(apis, l.svcCtx.DBEnt.API.Create().
		SetAPIGroup("Token").
		SetMethod("POST").
		SetPath("/token/unblock").
		SetServiceName("core").
		SetName("解封用户Token").SetIsRequired(false))

	err := l.svcCtx.DBEnt.API.CreateBulk(apis...).Exec(l.ctx)
	if err != nil {
		logx.Errorw("初始化基础API失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}

	return nil
}

// 初始化角色信息
func (l *InitLogic) initRole() error {
	var roles []*ent.RoleCreate

	roles = append(roles, l.svcCtx.DBEnt.Role.Create().
		SetRoleName("超级管理员").
		SetRoleCode("super").
		SetDescription("超级管理员 - 拥有所有权限"))

	err := l.svcCtx.DBEnt.Role.CreateBulk(roles...).Exec(l.ctx)
	if err != nil {
		logx.Errorw("初始化角色信息失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}

	return nil
}

// 初始化菜单信息
func (l *InitLogic) initMenu() error {
	var menus []*ent.MenuCreate

	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("dashboard").
		SetMenuName("仪表盘").
		SetMenuPath("/dashboard").
		SetComponent("/dashboard/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").SetIcon("carbon:dashboard"))

	// 系统管理（目录）
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("system").
		SetMenuName("系统管理").
		SetMenuPath("/system").
		SetSort(0).
		SetComponent("BasicLayout").
		SetMenuType("directory").
		SetServiceName("Core").
		SetIsHidden(false).
		SetIcon("carbon:settings"))

	// 系统管理 - 菜单
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("MenuManagement").
		SetMenuName("菜单管理").
		SetMenuPath("/system/menu").
		SetComponent("/system/menu/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("carbon:menu"))

	// 系统管理 - 角色
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("RoleManagement").
		SetMenuName("角色管理").
		SetMenuPath("/system/role").
		SetComponent("/system/role/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("carbon:user-role"))

	// 系统管理 - 部门
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("DepartmentManagement").
		SetMenuName("部门管理").
		SetMenuPath("/system/department").
		SetComponent("/system/department/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("carbon:user-multiple"))

	// 系统管理 - 用户
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("UserManagement").
		SetMenuName("用户管理").
		SetMenuPath("/system/user").
		SetComponent("/system/user/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("carbon:user"))

	// 系统管理 - 岗位
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("PositionManagement").
		SetMenuName("岗位管理").
		SetMenuPath("/system/position").
		SetComponent("/system/position/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("hugeicons:new-job"))

	// 系统管理 - 字典
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("DictManagement").
		SetMenuName("字典管理").
		SetMenuPath("/system/dict").
		SetComponent("/system/dict/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("material-symbols-light:dictionary"))
	// 系统管理 - API
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("ApiManagement").
		SetMenuName("API管理").
		SetMenuPath("/system/api").
		SetComponent("/system/api/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("carbon:api"))

	// 系统管理 - Oauth
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("OauthManagement").
		SetMenuName("Oauth管理").
		SetMenuPath("/system/oauth").
		SetComponent("/system/oauth/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("tabler:brand-oauth"))

	// 系统管理 - Token
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("TokenManagement").
		SetMenuName("Token管理").
		SetMenuPath("/system/token").
		SetComponent("/system/token/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("material-symbols:lock-open-outline"))

	// 系统管理 - 审计日志
	menus = append(menus, l.svcCtx.DBEnt.Menu.Create().
		SetMenuCode("AuditLogManagement").
		SetMenuName("审计日志").
		SetMenuPath("/system/audit").
		SetComponent("/system/audit/index").
		SetSort(0).
		SetServiceName("Core").
		SetMenuType("menu").
		SetParentID(2).
		SetIcon("tabler:info-square"))

	err := l.svcCtx.DBEnt.Menu.CreateBulk(menus...).Exec(l.ctx)
	if err != nil {
		logx.Errorw("初始化菜单信息失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}
	return nil
}

// 初始化OauthProvider信息
func (l *InitLogic) initOauthProvider() error {
	var providers []*ent.OauthProviderCreate

	// Google
	providers = append(providers, l.svcCtx.DBEnt.OauthProvider.Create().
		SetProviderName("Google").
		SetProviderCode("google").
		SetClientID("YOUR_CLIENT_ID").
		SetClientSecret("YOUR_CLIENT_SECRET").
		SetRedirectURI("http://localhost:8080/auth/oauth/callback").
		SetScopes("email openid").
		SetAuthorizationURL("https://accounts.google.com/o/oauth2/auth").
		SetTokenURL("https://oauth2.googleapis.com/token").
		SetUserinfoURL("https://www.googleapis.com/oauth2/v2/userinfo?access_token=TOKEN").
		SetAuthStyle(1).
		SetState(true))

	// Github
	providers = append(providers, l.svcCtx.DBEnt.OauthProvider.Create().
		SetProviderName("Github").
		SetProviderCode("github").
		SetClientID("YOUR_CLIENT_ID").
		SetClientSecret("YOUR_CLIENT_SECRET").
		SetRedirectURI("http://localhost:8080/auth/oauth/callback").
		SetScopes("email openid").
		SetAuthorizationURL("https://github.com/login/oauth/authorize").
		SetTokenURL("https://github.com/login/oauth/access_token").
		SetUserinfoURL("https://api.github.com/user").
		SetAuthStyle(2).
		SetState(true))

	err := l.svcCtx.DBEnt.OauthProvider.CreateBulk(providers...).Exec(l.ctx)
	if err != nil {
		logx.Errorw("初始化OauthProvider信息失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}
	return nil
}

// 初始化部门信息
func (l *InitLogic) initDepartment() error {
	var departments []*ent.DepartmentCreate

	departments = append(departments, l.svcCtx.DBEnt.Department.Create().
		SetDeptName("默认部门").
		SetDeptCode("default"))

	err := l.svcCtx.DBEnt.Department.CreateBulk(departments...).Exec(l.ctx)
	if err != nil {
		logx.Errorw("初始化部门信息失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}
	return nil
}

// 初始化职位信息
func (l *InitLogic) initPosition() error {
	var positions []*ent.PositionCreate
	positions = append(positions, l.svcCtx.DBEnt.Position.Create().
		SetPositionName("CEO").
		SetPositionCode("ceo").
		SetSort(0))

	err := l.svcCtx.DBEnt.Position.CreateBulk(positions...).Exec(l.ctx)
	if err != nil {
		logx.Errorw("初始化职位信息失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}
	return nil
}

// 初始化用户信息
func (l *InitLogic) initUser() error {
	var users []*ent.UserCreate
	users = append(users, l.svcCtx.DBEnt.User.Create().
		SetUsername("admin").
		SetPasswordHash(encrypt.BcryptEncrypt("last-admin123")).
		SetEmail("wenpiner@gmail.com").
		SetFullName("管理员").
		SetState(true).
		AddRoleIDs(1).
		SetDepartmentID(1).
		AddPositionIDs(1))

	err := l.svcCtx.DBEnt.User.CreateBulk(users...).Exec(l.ctx)
	if err != nil {
		logx.Errorw("初始化用户信息失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}
	return nil
}

// 初始化菜单和角色关联
func (l *InitLogic) initRoleMenu() error {
	// 获取当前所有菜单
	menus, err := l.svcCtx.DBEnt.Menu.Query().IDs(l.ctx)
	if err != nil {
		logx.Errorw("获取当前所有菜单失败", logx.Field("detail", err.Error()))
		return errorx.NewInternalError(err.Error())
	}
	// 更新角色
	err = l.svcCtx.DBEnt.Role.UpdateOneID(1).AddMenuIDs(menus...).Exec(l.ctx)
	return nil
}
