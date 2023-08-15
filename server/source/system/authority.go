package system

import (
	"context"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderAuthority = initOrderCasbin + 1

type initAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderAuthority, &initAuthority{})
}

func (i *initAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysAuthority{})
}

func (i *initAuthority) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysAuthority{})
}

func (i initAuthority) InitializerName() string {
	return sysModel.SysAuthority{}.TableName()
}

// 权限组(角色) 数据初始化
func (i *initAuthority) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []sysModel.SysAuthority{
		{AuthorityId: 888, AuthorityName: "超级管理员", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 1, AuthorityName: "普通玩家", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
	}
	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!", sysModel.SysAuthority{}.TableName())
	}

	// 超管关联角色 1和888   https://gorm.io/zh_CN/docs/associations.html#%E5%85%B3%E8%81%94%E6%A8%A1%E5%BC%8F
	if err := db.Model(&entities[0]).Association("DataAuthorityId").Replace(
		[]*sysModel.SysAuthority{
			{AuthorityId: 888},
			{AuthorityId: 1},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			db.Model(&entities[0]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}
	//普通玩家关联角色 1
	if err := db.Model(&entities[1]).Association("DataAuthorityId").Replace(
		[]*sysModel.SysAuthority{
			{AuthorityId: 1},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			db.Model(&entities[1]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}

	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("authority_id = ?", "1").
		First(&sysModel.SysAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
