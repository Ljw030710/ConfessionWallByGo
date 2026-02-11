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

// UpdatepasswordHandler API router注册点
func UpdatepasswordHandler() gin.HandlerFunc {
	api := UpdatepasswordApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfUpdatepassword).Pointer()).Name()] = api
	return hfUpdatepassword
}

type UpdatepasswordApi struct {
	Info     struct{}        `name:"修改密码" desc:"支持通过id或者用户名进行修改"`
	Request  UpdatepasswordApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response UpdatepasswordApiResponse // API响应数据 (Body中的Data部分)
}

type UpdatepasswordApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		UserID      int64  `json:"user_id" label:"用户ID"`
		Username    string `json:"username" label:"用户名"`
		NewPassword string `json:"new_password" binding:"required" label:"新密码"`
	}
}

type UpdatepasswordApiResponse struct {}

// Run Api业务逻辑执行点
func (u *UpdatepasswordApi) Run(ctx *gin.Context) kit.Code {
userRepo := repo.NewUserRepo()
	var err error

	// 1. 优先通过 UserID 修改
	if u.Request.Body.UserID > 0 {
		err = userRepo.UpdatePasswordByID(ctx, u.Request.Body.UserID, u.Request.Body.NewPassword)
	} else if u.Request.Body.Username != "" {
		// 2. 其次通过 Username 修改
		// 先查一下用户在不在
		user, _ := userRepo.FindByUsername(ctx, u.Request.Body.Username)
		if user == nil {
			return comm.CodeDataNotFound
		}
		err = userRepo.UpdatePasswordByUsername(ctx, u.Request.Body.Username, u.Request.Body.NewPassword)
	} else {
		return comm.CodeParameterInvalid
	}

	// 3. 处理数据库错误
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("数据库更新失败")
		return comm.CodeDatabaseError
	}

	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (u *UpdatepasswordApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindUri(&u.Request.Uri)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindHeader(&u.Request.Header)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindQuery(&u.Request.Query)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindJSON(&u.Request.Body)
	if err != nil {
		return err
	}
	return err
}

//  hfUpdatepassword API执行入口
func hfUpdatepassword(ctx *gin.Context) {
	api := &UpdatepasswordApi{}
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
