package api

import (
	"errors"
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
		UserID      int64  `json:"user_id" binding:"required" label:"用户ID"`
		OldPassword string `json:"old_password" binding:"required" label:"旧密码"`
		NewPassword string `json:"new_password" binding:"required" label:"新密码"`
	}
}

type UpdatepasswordApiResponse struct {}

// Run Api业务逻辑执行点
func (u *UpdatepasswordApi) Run(ctx *gin.Context) kit.Code {
userRepo := repo.NewUserRepo()

	// 1) 按 user_id 定位用户
	user, err := userRepo.FindById(ctx, u.Request.Body.UserID)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("查询用户失败")
		return comm.CodeDatabaseError
	}
	if user == nil {
		return comm.CodeDataNotFound
	}

	// 2) 校验旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Request.Body.OldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return comm.CodeAuthFailed
		}
		nlog.Pick().WithContext(ctx).WithError(err).Error("旧密码校验异常")
		return comm.CodePasswordEncryptError
	}

	// 3) 新密码加密后更新
	newHash, err := bcrypt.GenerateFromPassword([]byte(u.Request.Body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("新密码加密失败")
		return comm.CodePasswordEncryptError
	}

	err = userRepo.UpdatePasswordByID(ctx, u.Request.Body.UserID, string(newHash))
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("更新密码失败")
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
