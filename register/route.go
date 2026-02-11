package register

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/config"
	"github.com/zjutjh/mygo/middleware/cors"
	"github.com/zjutjh/mygo/swagger"

	"app/api"
)

func Route(router *gin.Engine) {
	router.Use(cors.Pick())

	r := router.Group(routePrefix())
	{
		routeBase(r, router)

		// 注册业务逻辑接口
		r.POST("/login",api.LoginHandler())//注册登录接口
		r.POST("/register",api.RegisterHandler())
		r.POST("/update_password",api.UpdatepasswordHandler())
		r.POST("/update_nickname",api.UpdatenicknameHandler())
		r.POST("/upload",api.UploadHandler())
		r.POST("/createconfession",api.CreateconfessionHandler())
		r.POST("/updateconfession",api.UpdateconfessionHandler())
		r.POST("/deleteconfession",api.DeleteconfessionHandler())
		r.POST("/block/user", api.BlockuserHandler())
		r.POST("/unblock/user", api.UnblockuserHandler())
		r.POST("/block/check", api.CheckblockHandler())
	}
}

func routePrefix() string {
	return "/api"
}

func routeBase(r *gin.RouterGroup, router *gin.Engine) {
	// OpenAPI/Swagger 文档生成
	if slices.Contains([]string{config.AppEnvDev, config.AppEnvTest}, config.AppEnv()) {
		r.GET("/swagger.json", swagger.DocumentHandler(router))
	}

	// 健康检查
	r.GET("/health", api.HealthHandler())
}
