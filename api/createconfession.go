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

// CreateconfessionHandler API router注册点
func CreateconfessionHandler() gin.HandlerFunc {
	api := CreateconfessionApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfCreateconfession).Pointer()).Name()] = api
	return hfCreateconfession
}

type CreateconfessionApi struct {
	Info     struct{}                    `name:"新增表白" desc:"创建一个新的表白记录"`
	Request  CreateconfessionApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response CreateconfessionApiResponse // API响应数据 (Body中的Data部分)
}

type CreateconfessionApiRequest struct {
	Uri    struct{}
	Header struct{}
	Query  struct{}
	Body   struct {
		SenderID     int64  `json:"sender_id" binding:"required" label:"发送者ID"`
		ReceiverName string `json:"receiver_name" binding:"required" label:"接收者昵称"`
		Content      string `json:"content" binding:"required" label:"表白内容"`
		ImageURL     string `json:"image_url" label:"图片链接"`
		IsAnonymous  int8   `json:"is_anonymous" label:"是否匿名(0/1)"`
		Status       int8   `json:"status" label:"状态(1/2/3)"`
	}
}

type CreateconfessionApiResponse struct {
	ID int64 `json:"id" desc:"新建表白"`
}

// Run Api业务逻辑执行点
func (c *CreateconfessionApi) Run(ctx *gin.Context) kit.Code {
	// TODO: 在此处编写接口业务逻辑
	confessionRepo := repo.NewConfessionRepo()
	//要先处理status
	status := c.Request.Body.Status
	if status == 0 {
		status = 1
	}
	newConfession, err := confessionRepo.Create(
		ctx,
		c.Request.Body.SenderID,
		c.Request.Body.ReceiverName,
		c.Request.Body.Content,
		c.Request.Body.ImageURL,
		c.Request.Body.IsAnonymous,
		status,
	)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("创建表白失败")
		return comm.CodeDatabaseError
	}
	c.Response.ID = newConfession.ID
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (c *CreateconfessionApi) Init(ctx *gin.Context) (err error) {
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

// hfCreateconfession API执行入口
func hfCreateconfession(ctx *gin.Context) {
	api := &CreateconfessionApi{}
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
