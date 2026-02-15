package register

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/config"
	"github.com/zjutjh/mygo/middleware/cors"
	"github.com/zjutjh/mygo/swagger"
	midjwt "github.com/zjutjh/mygo/jwt/middleware"
	"app/api"
)

func Route(router *gin.Engine) {
	router.Use(cors.Pick())

	r := router.Group(routePrefix())
	{
		routeBase(r, router)
		auth := midjwt.Auth[int64](true)
		// 注册业务逻辑接口
		r.POST("/user/login",api.LoginHandler())//注册登录接口
		r.POST("/user/register",api.RegisterHandler())
		r.POST("/user/update_password", auth, api.UpdatepasswordHandler())
		r.POST("/user/update_nickname", auth, api.UpdatenicknameHandler())
		r.POST("/confession/upload", auth, api.UploadHandler())
		r.POST("/confession/createconfession", auth, api.CreateconfessionHandler())
		r.POST("/confession/updateconfession", auth, api.UpdateconfessionHandler())
		r.POST("/confession/deleteconfession", auth, api.DeleteconfessionHandler())
		r.POST("/block/user", auth, api.BlockuserHandler())
		r.POST("/unblock/user", auth, api.UnblockuserHandler())
		r.POST("/block/check", auth, api.CheckblockHandler())
		r.POST("/comment/createcomment", auth, api.CreatecommentHandler())
		r.POST("/comment/replycomment", auth, api.ReplycommentHandler())
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
