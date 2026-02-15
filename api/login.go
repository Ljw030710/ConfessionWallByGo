package api

import (
	"errors"
	"reflect"
	"runtime"
	midjwt "github.com/zjutjh/mygo/jwt"
	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"
	"golang.org/x/crypto/bcrypt"

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
	ID       int64  `json:"id" desc:"用户ID"`
	Username string `json:"username" desc:"用户名"`
	Nickname string `json:"nickname" desc:"昵称"`
	Token    string `json:"token" desc:"JWT Token"`
}

// Run Api业务逻辑执行点
func (l *LoginApi) Run(ctx *gin.Context) kit.Code {
	//根据用户名查询用户
	userRepo := repo.NewUserRepo()
	user,err := userRepo.FindByUsername(ctx,l.Request.Body.Username)
	if err != nil{
		nlog.Pick().WithContext(ctx).WithError(err).Error("用户查询失败")
		return comm.CodeDatabaseError
	}
	if user == nil{
		return comm.CodeAuthFailed
	}
	//进行密码校验
	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(l.Request.Body.Password))
	if err != nil{
		if errors.Is(err,bcrypt.ErrMismatchedHashAndPassword){
			return comm.CodeAuthFailed
		}
		nlog.Pick().WithContext(ctx).WithError(err).Error("密码校验异常")
		return comm.CodeAuthFailed
	}
	//密码正确发JWT
	token,err := midjwt.Pick[int64]().GenerateToken(user.ID)
	if err != nil{
		nlog.Pick().WithContext(ctx).WithError(err).Error("生成JWT失败")
		return comm.CodeMiddlewareServiceError
	}
	l.Response.ID = user.ID
	l.Response.Username = user.Username
	l.Response.Nickname = user.Nickname
	l.Response.Token = token
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
