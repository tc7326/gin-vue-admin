package register

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/router"
	"github.com/gin-gonic/gin"
)

type RegisterPlugin struct {
}

func CreateRegisterPlug(AuthorityId uint) *RegisterPlugin {
	//注册插件的同时 设置 玩家注册的默认用户组
	global.GlobalConfig.AuthorityId = AuthorityId
	return &RegisterPlugin{}
}

func (*RegisterPlugin) Register(group *gin.RouterGroup) {
	router.RouterGroupApp.InitRegisterRouter(group)
}

// RouterPath 当前插件的base路由的路径
func (*RegisterPlugin) RouterPath() string {
	return "register"
}
