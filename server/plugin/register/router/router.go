package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/api"
	"github.com/gin-gonic/gin"
)

type RegisterRouter struct {
}

func (s *RegisterRouter) InitRegisterRouter(Router *gin.RouterGroup) {
	plugRouter := Router
	plugApi := api.ApiGroupApp.RegisterApi
	{
		//插件案例的 直接注册的接口
		plugRouter.POST("register", plugApi.ApiName)

		//用户通过邮箱验证自己注册的注册接口
		plugRouter.POST("user", plugApi.UserRegister)
	}
}
