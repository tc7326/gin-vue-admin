package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		baseRouter.POST("login", baseApi.Login)            //登录
		baseRouter.POST("captcha", baseApi.Captcha)        //登录验证码
		baseRouter.POST("captcha-e", baseApi.CaptchaEmail) //发送邮箱验证码
	}
	return baseRouter
}
