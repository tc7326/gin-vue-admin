package api

import (
	systemApi "github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/captcha"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

var store = captcha.NewDefaultRedisStore() //使用redis缓存

type RegisterApi struct{}

// ApiName
// @Tags Register
// @Summary 请手动填写接口功能
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /register/routerName[post]
func (p *RegisterApi) ApiName(c *gin.Context) {
	var plug model.Request
	_ = c.ShouldBindJSON(&plug)
	if res, err := service.ServiceGroupApp.PlugService(plug); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		var baseApi systemApi.BaseApi
		baseApi.TokenNext(c, *res)
	}
}

// UserRegister
// @Tags Register
// @Summary 用户通过邮箱验证码自己注册的注册接口
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /register/user[post]
func (p *RegisterApi) UserRegister(c *gin.Context) {
	//结构体 校验
	var req model.Request
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	log.Printf("用户注册传递的参数是: %v", req)
	//参数空校验
	err = utils.Verify(req, utils.EmailRegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//校验验证码
	if req.CaptchaId == "" {
		response.FailWithMessage("请先获取验证码", c)
		return
	}

	//redis使用ctx
	store.UseWithCtx(c)
	//获取缓存的验证码
	captchaCode := store.VerifyGet(req.CaptchaId)
	log.Printf("缓存读取到的验证码是: %v", captchaCode)
	if captchaCode == "" {
		response.FailWithMessage("验证码已过期 请重新获取", c)
		return
	}
	if req.Captcha != captchaCode {
		response.FailWithMessage("验证码错误 请重试", c)
		return
	}
	//获取缓存的用户名
	capUsername := store.CapInfoGet(req.CaptchaId)
	log.Printf("缓存读取到的验证码对应的用户名是: %v", capUsername)
	if capUsername == "" {
		response.FailWithMessage("验证码已过期 请重新获取!", c)
		return
	}
	if req.Username != capUsername {
		response.FailWithMessage("验证码错误 请重试!", c)
		return
	}

	//设置关联权限组 这里只有1 普通用户
	authorities := make([]system.SysAuthority, 1)
	authorities[0] = system.SysAuthority{AuthorityId: 1}

	//默认头像
	steve := "https://gss0.baidu.com/-fo3dSag_xI4khGko9WTAnF6hhy/zhidao/wh%3D600%2C800/sign=0039b765d2160924dc70aa1de43719c2/bd315c6034a85edf3b752e104b540923dd54750c.jpg"

	//组装业务数据
	user := &system.SysUser{
		Username:    req.Username,
		NickName:    req.Username,
		Password:    req.Password,
		AuthorityId: 1,           //用户当前的用户组
		Authorities: authorities, //用户所有的用户组
		HeaderImg:   steve,       //用户头像 默认史蒂夫
		Email:       req.Email,
	}

	//具体业务
	userReturn, err := service.ServiceGroupApp.Register(*user)
	if err != nil {
		global.GVA_LOG.Error("注册失败!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册失败 "+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册成功", c)
}
