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

// UpdatenicknameHandler API router注册点
func UpdatenicknameHandler() gin.HandlerFunc {
	api := UpdatenicknameApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfUpdatenickname).Pointer()).Name()] = api
	return hfUpdatenickname
}

type UpdatenicknameApi struct {
	Info     struct{}        `name:"修改昵称" desc:"根据用户名进行修改"`
	Request  UpdatenicknameApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response UpdatenicknameApiResponse // API响应数据 (Body中的Data部分)
}

type UpdatenicknameApiRequest struct {
	Uri struct {}
	Header struct {}
	Query struct {}
	Body struct {
		Username string `json:"username" binding:"required" label:"用户名"`
		NewNickname string `json:"new_nickname" binding:"required" label:"新昵称"`
	}
}

type UpdatenicknameApiResponse struct {}

func (u *UpdatenicknameApi) Run(ctx *gin.Context) kit.Code {
	userRepo := repo.NewUserRepo()

	// 1. 直接通过用户名修改昵称
	err := userRepo.UpdateNicknameByUsername(ctx, u.Request.Body.Username, u.Request.Body.NewNickname)
	
	// 2. 错误处理
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comm.CodeDataNotFound // 如果用户名不存在
		}
		nlog.Pick().WithContext(ctx).WithError(err).Error("更新昵称失败")
		return comm.CodeDatabaseError
	}

	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (u *UpdatenicknameApi) Init(ctx *gin.Context) (err error) {
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

//  hfUpdatenickname API执行入口
func hfUpdatenickname(ctx *gin.Context) {
	api := &UpdatenicknameApi{}
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
