package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenuAuthority = initOrderMenu + initOrderAuthority

type initMenuAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initMenuAuthority) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i initMenuAuthority) InitializerName() string {
	return "sys_menu_authorities"
}

// 菜单和权限对应关系 数据初始化
func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	//权限表 这是拿到所有权限了?
	authorities, ok := ctx.Value(initAuthority{}.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}

	//菜单表 这是拿到所有菜单了?
	menus, ok := ctx.Value(initMenu{}.InitializerName()).([]sysModel.SysBaseMenu)
	if !ok {
		return next, errors.Wrap(errors.New(""), "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx

	// 将 所有按钮 的权限 关联给 超管角色(888) 登录后就可能看到什么按钮(选项卡)
	if err = db.Model(&authorities[0]).Association("SysBaseMenus").Replace(menus); err != nil {
		return next, err
	}

	// 将 仪表盘 和 关于 按钮权限 关联给 普通玩家(1) 按钮权限太多 我看不过来 这里就给2个基础 后续在数据库添加
	menu1 := menus[:2]
	if err = db.Model(&authorities[1]).Association("SysBaseMenus").Replace(menu1); err != nil {
		return next, err
	}
	return next, nil
}

func (i *initMenuAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	auth := &sysModel.SysAuthority{}
	if ret := db.Model(auth).
		Where("authority_id = ?", 1).Preload("SysBaseMenus").Find(auth); ret != nil {
		if ret.Error != nil {
			return false
		}
		return len(auth.SysBaseMenus) > 0
	}
	return false
}
