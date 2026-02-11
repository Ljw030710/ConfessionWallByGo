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

// UpdateconfessionHandler API router注册点
func UpdateconfessionHandler() gin.HandlerFunc {
	api := UpdateconfessionApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfUpdateconfession).Pointer()).Name()] = api
	return hfUpdateconfession
}

type UpdateconfessionApi struct {
	Info     struct{}        `name:"修改表白" desc:"根据confession_id 修改表白"`
	Request  UpdateconfessionApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response UpdateconfessionApiResponse // API响应数据 (Body中的Data部分)
}

type UpdateconfessionApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		ConfessionID int64 `json:"confession_id" binding:"required" label:"表白ID"`
		// 以下为可选字段：传了就更新，不传就保持原值
		ReceiverName *string `json:"receiver_name" label:"接收者昵称"`
		Content      *string `json:"content" label:"表白内容"`
		ImageURL     *string `json:"image_url" label:"图片链接"`
		IsAnonymous  *int8   `json:"is_anonymous" label:"是否匿名(0/1)"`
		Status       *int8   `json:"status" label:"状态(1/2/3)"`
	}
}

type UpdateconfessionApiResponse struct {}

// Run Api业务逻辑执行点
func (u *UpdateconfessionApi) Run(ctx *gin.Context) kit.Code {
	// TODO: 在此处编写接口业务逻辑
	//1、至少需要一个待更新字段
	confessionRepo := repo.NewConfessionRepo()
	if (u.Request.Body.ReceiverName == nil &&u.Request.Body.Content == nil&&u.Request.Body.ImageURL == nil &&u.Request.Body.IsAnonymous == nil &&u.Request.Body.Status == nil) {
		return comm.CodeParameterInvalid
	}
	err := confessionRepo.UpdateByID(
		ctx,
		u.Request.Body.ConfessionID,
		u.Request.Body.ReceiverName,
		u.Request.Body.Content,
		u.Request.Body.ImageURL,
		u.Request.Body.IsAnonymous,
		u.Request.Body.Status,
	)

	if err != nil {
		// 4) 区分“数据不存在”和“数据库异常”
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comm.CodeDataNotFound
		}
		nlog.Pick().WithContext(ctx).WithError(err).Error("更新表白失败")
		return comm.CodeDatabaseError
	}

	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (u *UpdateconfessionApi) Init(ctx *gin.Context) (err error) {
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

//  hfUpdateconfession API执行入口
func hfUpdateconfession(ctx *gin.Context) {
	api := &UpdateconfessionApi{}
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
