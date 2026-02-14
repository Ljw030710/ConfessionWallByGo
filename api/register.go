package api

import (
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"
	"golang.org/x/crypto/bcrypt"

	"app/comm"
	"app/dao/repo"
)

// RegisterHandler API router注册点
func RegisterHandler() gin.HandlerFunc {
	api := RegisterApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfRegister).Pointer()).Name()] = api
	return hfRegister
}

type RegisterApi struct {
	Info     struct{}        `name:"用户注册" desc:"创建新账号"`
	Request  RegisterApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response RegisterApiResponse // API响应数据 (Body中的Data部分)
}

type RegisterApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		Username string `json:"username" binding:"required" label:"用户名"`
		Password string `json:"password" binding:"required" label:"密码"`
		Nickname string `json:"nickname" binding:"required" label:"昵称"`

	}
}

type RegisterApiResponse struct {
	ID int64 `json:"id" desc:"新用户ID"`
}

// Run Api业务逻辑执行点
func (r *RegisterApi) Run(ctx *gin.Context) kit.Code {
	// TODO: 在此处编写接口业务逻辑
	// 1. 初始化 Repo
	userRepo := repo.NewUserRepo()

	// 2) 使用 bcrypt 对明文密码做单向哈希
	// 注意：哈希后不可逆，登录时用 CompareHashAndPassword 校验
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Request.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("密码加密失败")
		return comm.CodePasswordEncryptError
	}

	// 3) 入库时只保存哈希后的密码，绝不保存明文
	newUser, err := userRepo.Create(ctx, r.Request.Body.Username, string(hashedPassword), r.Request.Body.Nickname)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("创建用户失败")
		return comm.CodeDatabaseError
	}

	// 4) 返回新用户ID
	r.Response.ID = newUser.ID
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (r *RegisterApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindUri(&r.Request.Uri)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindHeader(&r.Request.Header)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindQuery(&r.Request.Query)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindJSON(&r.Request.Body)
	if err != nil {
		return err
	}
	return err
}

//  hfRegister API执行入口
func hfRegister(ctx *gin.Context) {
	api := &RegisterApi{}
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
