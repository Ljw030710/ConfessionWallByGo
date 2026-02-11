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
	"gorm.io/gorm"

	"app/comm"
	"app/dao/repo"
)

// UnblockuserHandler API router注册点
func UnblockuserHandler() gin.HandlerFunc {
	api := UnblockuserApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfUnblockuser).Pointer()).Name()] = api
	return hfUnblockuser
}

type UnblockuserApi struct {
	Info     struct{}        `name:"取消拉黑用户" desc:"通过用户名取消拉黑"`
	Request  UnblockuserApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response UnblockuserApiResponse // API响应数据 (Body中的Data部分)
}

type UnblockuserApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		BlockerUsername string `json:"blocker_username" binding:"required" label:"拉黑发起人用户名"`
		BlockedUsername string `json:"blocked_username" binding:"required" label:"被拉黑用户名"`
	}
}

type UnblockuserApiResponse struct {}

// Run Api业务逻辑执行点
func (u *UnblockuserApi) Run(ctx *gin.Context) kit.Code {
	// TODO: 在此处编写接口业务逻辑
	userBlockRepo := repo.NewUserBlockRepo()
	err := userBlockRepo.UnblockByUsername(ctx, u.Request.Body.BlockerUsername, u.Request.Body.BlockedUsername)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comm.CodeDataNotFound
		}
		nlog.Pick().WithContext(ctx).WithError(err).Error("取消拉黑失败")
		return comm.CodeDatabaseError
	}
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (u *UnblockuserApi) Init(ctx *gin.Context) (err error) {
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

//  hfUnblockuser API执行入口
func hfUnblockuser(ctx *gin.Context) {
	api := &UnblockuserApi{}
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
