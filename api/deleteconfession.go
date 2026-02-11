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

// DeleteconfessionHandler API router注册点
func DeleteconfessionHandler() gin.HandlerFunc {
	api := DeleteconfessionApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfDeleteconfession).Pointer()).Name()] = api
	return hfDeleteconfession
}

type DeleteconfessionApi struct {
	Info     struct{}        `name:"删除表白" desc:"根据confession_id来删除表白"`
	Request  DeleteconfessionApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response DeleteconfessionApiResponse // API响应数据 (Body中的Data部分)
}

type DeleteconfessionApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		ConfessionID int64 `json:"confession_id" binding:"required" label:"表白ID"`
	}
}

type DeleteconfessionApiResponse struct {}

// Run Api业务逻辑执行点
func (d *DeleteconfessionApi) Run(ctx *gin.Context) kit.Code {
	// 1) 初始化仓储对象
	confessionRepo := repo.NewConfessionRepo()

	// 2) 执行删除（建议 repo 内部做逻辑删除：status = 3）
	err := confessionRepo.DeleteByID(ctx, d.Request.Body.ConfessionID)
	if err != nil {
		// 3) 记录不存在：返回业务“数据不存在”
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comm.CodeDataNotFound
		}
		// 4) 其他错误：按数据库错误处理并记录日志
		nlog.Pick().WithContext(ctx).WithError(err).Error("删除表白失败")
		return comm.CodeDatabaseError
	}

	// 5) 删除成功
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (d *DeleteconfessionApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindUri(&d.Request.Uri)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindHeader(&d.Request.Header)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindQuery(&d.Request.Query)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindJSON(&d.Request.Body)
	if err != nil {
		return err
	}
	return err
}

//  hfDeleteconfession API执行入口
func hfDeleteconfession(ctx *gin.Context) {
	api := &DeleteconfessionApi{}
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
