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

// CreatecommentHandler API router注册点
func CreatecommentHandler() gin.HandlerFunc {
	api := CreatecommentApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfCreatecomment).Pointer()).Name()] = api
	return hfCreatecomment
}

type CreatecommentApi struct {
	Info     struct{}        `name:"新增评论" desc:"给指定表白新增一条评论"`
	Request  CreatecommentApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response CreatecommentApiResponse // API响应数据 (Body中的Data部分)
}

type CreatecommentApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		ConfessionID int64 `json:"confession_id" binding:"required" label:"表白ID"`
		Username string `json:"username" binding:"required" label:"评论用户名"`
		Content string `json:"content" binding:"required" label:"评论内容"`
	}
}

type CreatecommentApiResponse struct {
	ID int64 `json:"confession_id" binding:"required" label:"新键评论id"`
}

// Run Api业务逻辑执行点
func (c *CreatecommentApi) Run(ctx *gin.Context) kit.Code {
	// TODO: 在此处编写接口业务逻辑
	commentRepo := repo.NewConfessionCommentRepo()
	newComment,err := commentRepo.Create(
		ctx,
		c.Request.Body.ConfessionID,
		c.Request.Body.Username,
		c.Request.Body.Content,
	)
	if err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return comm.CodeDataNotFound
		}
		nlog.Pick().WithContext(ctx).WithError(err).Error("创建评论失败")
		return comm.CodeDatabaseError
	}
	c.Response.ID = newComment.ID
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (c *CreatecommentApi) Init(ctx *gin.Context) (err error) {
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

//  hfCreatecomment API执行入口
func hfCreatecomment(ctx *gin.Context) {
	api := &CreatecommentApi{}
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
