package system

import (
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	email_response "github.com/flipped-aurora/gin-vue-admin/server/plugin/email/model/response"
	emailService "github.com/flipped-aurora/gin-vue-admin/server/plugin/email/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/captcha"
	"log"
	"os"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
var store = captcha.NewDefaultRedisStore() //使用redis缓存
//var store = base64Captcha.DefaultMemStore // 使用内存

type BaseApi struct{}

// Captcha
// @Tags      Base
// @Summary   生成验证码
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysCaptchaResponse,msg=string}  "生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router    /base/captcha [post]
func (b *BaseApi) Captcha(c *gin.Context) {
	// 判断验证码是否开启
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	key := c.ClientIP()
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool
	if openCaptcha == 0 || openCaptcha < interfaceToInt(v) {
		oc = true
	}
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.GVA_CONFIG.Captcha.ImgHeight, global.GVA_CONFIG.Captcha.ImgWidth, global.GVA_CONFIG.Captcha.KeyLong, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c)) // v8下使用redis
	//cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		global.GVA_LOG.Error("验证码获取失败!", zap.Error(err))
		response.FailWithMessage("验证码获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysCaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: global.GVA_CONFIG.Captcha.KeyLong,
		OpenCaptcha:   oc,
	}, "验证码获取成功", c)
}

// 类型转换
func interfaceToInt(v interface{}) (i int) {
	switch v := v.(type) {
	case int:
		i = v
	default:
		i = 0
	}
	return
}

// CaptchaEmail 邮箱验证码 生成验证码的同时 下发邮件
func (b *BaseApi) CaptchaEmail(c *gin.Context) {

	//复用注册的结构体 主要用于获取用户邮箱
	var l systemReq.Login
	err := c.ShouldBindJSON(&l)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//校验邮箱格式
	if len(l.Email) < 8 || !strings.Contains(l.Email, "@qq.com") {
		response.FailWithMessage("邮箱格式不正确", c)
		return
	}

	// 判断验证码是否开启
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	key := c.ClientIP()
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool
	if openCaptcha == 0 || openCaptcha < interfaceToInt(v) {
		oc = true
	}

	// 这个应该是生成验证码的驱动
	driver := base64Captcha.NewDriverDigit(global.GVA_CONFIG.Captcha.ImgHeight, global.GVA_CONFIG.Captcha.ImgWidth, global.GVA_CONFIG.Captcha.KeyLong, 0.7, 80)
	// 这里应该是生成了验证码的图的base64
	cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c)) // v8下使用redis
	//cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		global.GVA_LOG.Error("验证码获取失败!", zap.Error(err))
		response.FailWithMessage("验证码获取失败", c)
		return
	}
	log.Println("邮箱验证码缓存的key: ", id)

	//然后读取缓存 拿到真的验证码(由于验证码生成逻辑封装了 还没看懂 所以我直接从缓存里读)
	captchaCode := store.VerifyGet(id)

	log.Println("邮箱验证码是: ", captchaCode)

	//这里加下发邮件的逻辑

	//读取邮件模板
	tempBytes, err := os.ReadFile("./plugin/email/template/user_reg.html")
	if err != nil {
		global.GVA_LOG.Error("模板读取失败!", zap.Error(err))
		response.FailWithMessage("模板读取失败", c)
		return
	}

	//替换内容
	html := strings.Replace(string(tempBytes), "123456", captchaCode, 1)
	//log.Println("替换后的html", html)

	//这里直接套用 下发邮件的结构体
	var email email_response.Email
	email.To = l.Email           //用户注册用的邮箱
	email.Subject = "红土大陆注册安全验证" //标题
	email.Body = html            //填充html模板

	err = emailService.ServiceGroupApp.SendEmail(email.To, email.Subject, email.Body)
	if err != nil {
		global.GVA_LOG.Error("发送失败!", zap.Error(err))
		response.FailWithMessage("邮件发送失败 请重试", c)
		return
	}

	//返回成功 需要将 CaptchaId(会校验ip地址) 返回给前端 注册时需要校验
	response.OkWithDetailed(systemRes.SysCaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: global.GVA_CONFIG.Captcha.KeyLong,
		OpenCaptcha:   oc,
	}, "验证邮件已发送 有效期5分钟", c)
}
