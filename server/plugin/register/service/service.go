package service

import (
	"errors"
	"github.com/gofrs/uuid/v5"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	plugGlobal "github.com/flipped-aurora/gin-vue-admin/server/plugin/register/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/model"
	userService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/mojocn/base64Captcha"
)

type RegisterService struct{}

func (e *RegisterService) PlugService(req model.Request) (res *system.SysUser, err error) {

	//空值校验
	if err := utils.Verify(req, utils.LoginVerify); err != nil {
		// 某个参数为空
		return res, err
	}
	var (
		store = base64Captcha.DefaultMemStore
		user  system.SysUser
		us    *userService.UserService
	)
	if !store.Verify(req.CaptchaId, req.Captcha, true) {
		return res, errors.New("验证码错误")
	}

	//先查数据库有没有这个人
	u := &system.SysUser{Username: req.Username, Password: req.Password}
	err = global.GVA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		return res, errors.New("用户名已注册")
	}
	if user.Username != "" && user.Password != "" {
		return res, errors.New("用户名已注册")
	}

	var sysAuthority systemReq.Register
	sysAuthority.Username = u.Username
	sysAuthority.NickName = u.NickName
	sysAuthority.Password = u.Password
	//设置权限组
	sysAuthority.AuthorityId = plugGlobal.GlobalConfig.AuthorityId
	sysAuthority.AuthorityIds = append(sysAuthority.AuthorityIds, plugGlobal.GlobalConfig.AuthorityId)

	// 定义 新注册的用户的结构体
	user.Password = u.Password
	user.UUID = uuid.Must(uuid.NewV4())
	user.Username = u.Username
	user.NickName = u.Username
	user.AuthorityId = plugGlobal.GlobalConfig.AuthorityId

	for _, v := range sysAuthority.AuthorityIds {
		user.Authorities = append(user.Authorities, system.SysAuthority{
			AuthorityId:   v,
			DefaultRouter: "dashboard", //配置默认路由 就用户登录后默认是那个页面
		})
	}

	if rest, err := us.Register(*u); err != nil {
		return &rest, errors.New("注册失败!")
	}
	if res, err = us.Login(u); err != nil {
		return res, errors.New("登陆失败!")
	}
	return res, nil
	// 前面的代码 拿不到正确的 user，所以需要再次查询一次

}
