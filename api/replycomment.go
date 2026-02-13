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

// ReplycommentHandler API router注册点
func ReplycommentHandler() gin.HandlerFunc {
	api := ReplycommentApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfReplycomment).Pointer()).Name()] = api
	return hfReplycomment
}

type ReplycommentApi struct {
	Info     struct{}        `name:"回复评论" desc:"回复某条评论"`
	Request  ReplycommentApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response ReplycommentApiResponse // API响应数据 (Body中的Data部分)
}

type ReplycommentApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		ConfessionID int64 `json:"confession_id" binding:"required" label:"表白ID"`
		ParentCommentID int64  `json:"parent_comment_id" binding:"required" label:"父评论ID"`
		Username        string `json:"username" binding:"required" label:"回复人用户名"`
		ReplyToUsername string `json:"reply_to_username" label:"被回复用户名(可选)"`
		Content         string `json:"content" binding:"required" label:"回复内容"`
	}
}

type ReplycommentApiResponse struct {
	ID int64 `json:"id" desc:"回复ID"`
}

// Run Api业务逻辑执行点
func (r *ReplycommentApi) Run(ctx *gin.Context) kit.Code {
	// TODO: 在此处编写接口业务逻辑
	commentRepo := repo.NewConfessionCommentRepo()
	//如果ID《0给个非法
	if r.Request.Body.ParentCommentID<=0{
		return comm.CodeParameterInvalid
	}
	//调用repo的reply接口函数
	newReply,err := commentRepo.Reply(
		ctx,
		r.Request.Body.ConfessionID,
		r.Request.Body.ParentCommentID,
		r.Request.Body.Username,
		r.Request.Body.ReplyToUsername,
		r.Request.Body.Content,
	)
	if err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return comm.CodeDataNotFound
		}
		nlog.Pick().WithContext(ctx).WithError(err).Error("回复评论失败")
		return comm.CodeDatabaseError
	}
	r.Response.ID = newReply.ID
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (r *ReplycommentApi) Init(ctx *gin.Context) (err error) {

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

//  hfReplycomment API执行入口
func hfReplycomment(ctx *gin.Context) {
	api := &ReplycommentApi{}
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
