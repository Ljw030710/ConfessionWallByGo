package api

import (
	"errors"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"
	"gorm.io/gorm"

	"app/comm"
	"app/dao/repo"
)

// BlockuserHandler API router注册点
func BlockuserHandler() gin.HandlerFunc {
	api := BlockuserApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfBlockuser).Pointer()).Name()] = api
	return hfBlockuser
}

type BlockuserApi struct {
	Info     struct{}        `name:"拉黑用户" desc:"用用户名拉黑用户"`
	Request  BlockuserApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response BlockuserApiResponse // API响应数据 (Body中的Data部分)
}

type BlockuserApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		BlockerUsername string `json:"blocker_username" binding:"required" label:"拉黑发起人用户名"`
		BlockedUsername string `json:"blocked_username" binding:"required" label:"被拉黑用户名"`
	}
}

type BlockuserApiResponse struct {}

// Run Api业务逻辑执行点
func (b *BlockuserApi) Run(ctx *gin.Context) kit.Code {
	userBlockRepo := repo.NewUserBlockRepo()
	err := userBlockRepo.BlockByUsername(ctx, b.Request.Body.BlockerUsername, b.Request.Body.BlockedUsername)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comm.CodeDataNotFound
		}
		// repo里“自己拉黑自己”是errors.New("不能拉黑自己")
		if strings.Contains(err.Error(), "不能拉黑自己") {
			return comm.CodeParameterInvalid
		}
		// 唯一索引冲突：重复拉黑
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return comm.CodeDataConflict
		}
		nlog.Pick().WithContext(ctx).WithError(err).Error("拉黑用户失败")
		return comm.CodeDatabaseError
	}
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (b *BlockuserApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindUri(&b.Request.Uri)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindHeader(&b.Request.Header)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindQuery(&b.Request.Query)
	if err != nil {
		return err
	}
	err = ctx.ShouldBindJSON(&b.Request.Body)
	if err != nil {
		return err
	}
	return err
}

//  hfBlockuser API执行入口
func hfBlockuser(ctx *gin.Context) {
	api := &BlockuserApi{}
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
