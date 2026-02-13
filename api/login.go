package api

import (
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"

	"app/comm"
	"app/dao/repo"
)

// LoginHandler API router注册点
func LoginHandler() gin.HandlerFunc {
	api := LoginApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfLogin).Pointer()).Name()] = api
	return hfLogin
}

type LoginApi struct {
	Info     struct{}        `name:"用户登录" desc:"用户名和密码登录"`
	Request  LoginApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response LoginApiResponse // API响应数据 (Body中的Data部分)
}

type LoginApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		Username string `json:"username" binding:"required" label:"用户名"`
		Password string `json:"password" binding:"required" label:"密码"`
	}
}

type LoginApiResponse struct {
	ID int64 `json:"id" desc:"用户ID"`
	Username string `json:"username" desc:"用户名"`
	Nickname string `json:"nickname" desc:"昵称"`
}

// Run Api业务逻辑执行点
func (l *LoginApi) Run(ctx *gin.Context) kit.Code {
	// TODO: 在此处编写接口业务逻辑
	//1、初始化
	userRepo := repo.NewUserRepo()
	//2、执行登录逻辑
	user,err := userRepo.Login(ctx,l.Request.Body.Username,l.Request.Body.Password)
	if err != nil{
		nlog.Pick().WithContext(ctx).WithError(err).Warn("登录验证失败")
		if err.Error() == "用户名或密码错误"{
			return comm.CodeAuthFailed
		}
		return comm.CodeDatabaseError
	}

	//3、填充响应数据
	l.Response.ID = user.ID
	l.Response.Username  =user.Username
	l.Response.Nickname = user.Nickname
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (l *LoginApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindUri(&l.Request.Uri)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindHeader(&l.Request.Header)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindQuery(&l.Request.Query)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindJSON(&l.Request.Body)
	if err != nil {
		return err
	}
	return err
}

//  hfLogin API执行入口
func hfLogin(ctx *gin.Context) {
	api := &LoginApi{}
	err := api.Init(ctx)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Warn("参数绑定校验错误")
		reply.Fail(ctx, comm.CodeParameterInvalid)
		return
	}
	code := api.Run(ctx)
	if !ctx.IsAborted() {
		if code == comm.CodeOK {
			reply.Success(ctx, api.Response)
		} else {
			reply.Fail(ctx, code)
		}
	}
}
