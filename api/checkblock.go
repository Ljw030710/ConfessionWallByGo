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

// CheckblockHandler API router注册点
func CheckblockHandler() gin.HandlerFunc {
	api := CheckblockApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfCheckblock).Pointer()).Name()] = api
	return hfCheckblock
}

type CheckblockApi struct {
	Info     struct{}        `name:"查询是否拉黑" desc:"通过用户名来判断是否拉黑"`
	Request  CheckblockApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response CheckblockApiResponse // API响应数据 (Body中的Data部分)
}

type CheckblockApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		BlockerUsername string `json:"blocker_username" binding:"required" label:"拉黑发起人用户名"`
		BlockedUsername string `json:"blocked_username" binding:"required" label:"被拉黑用户名"`
	}
}

type CheckblockApiResponse struct {
	IsBlocked bool `json:"is_blocked"` // true=已拉黑，false=未拉黑
}

// Run Api业务逻辑执行点
func (c *CheckblockApi) Run(ctx *gin.Context) kit.Code {
	userBlockRepo := repo.NewUserBlockRepo()
	isBlocked, err := userBlockRepo.IsBlockedByUsername(ctx, c.Request.Body.BlockerUsername, c.Request.Body.BlockedUsername)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comm.CodeDataNotFound
		}
		nlog.Pick().WithContext(ctx).WithError(err).Error("查询拉黑关系失败")
		return comm.CodeDatabaseError
	}
	c.Response.IsBlocked = isBlocked
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (c *CheckblockApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindUri(&c.Request.Uri)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindHeader(&c.Request.Header)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindQuery(&c.Request.Query)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindJSON(&c.Request.Body)
	if err != nil {
		return err
	}
	return err
}

//  hfCheckblock API执行入口
func hfCheckblock(ctx *gin.Context) {
	api := &CheckblockApi{}
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
