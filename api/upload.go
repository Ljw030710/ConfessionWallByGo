package api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"

	"app/comm"
	"app/dao/repo"
)

const (
	qiniuAccessKey = "7Esv5lWwgMhEoCUnuJbCY-oDTnijXgRzll5oNCZh"
	qiniuSecretKey = "QdLH0zgtZyMt8bCxSPdFXT4zLtwYEiJkTZLjv4Fa"
	qiniuBucket    = "repoaa"
	qiniuDomain    = "http://ta4qfny94.hd-bkt.clouddn.com/"
)

// UploadHandler API router注册点
func UploadHandler() gin.HandlerFunc {
	api := UploadApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfUpload).Pointer()).Name()] = api
	return hfUpload
}

type UploadApi struct {
	Info     struct{}          `name:"图片上传" desc:"支持MD5去重并上传到七牛"`
	Request  UploadApiRequest  // API请求参数 (Uri/Header/Query/Form)
	Response UploadApiResponse // API响应数据
}

type UploadApiRequest struct {
	Uri    struct{}
	Header struct{}
	Query  struct{}
	Form   struct {
		Type string `form:"type" json:"type"`
	}
}

type UploadApiResponse struct {
	URL  string `json:"url"`
	Hash string `json:"hash"`
}

// Run Api业务逻辑执行点
func (u *UploadApi) Run(ctx *gin.Context) kit.Code {
	
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return comm.CodeParameterInvalid
	}

	if fileHeader.Size > 5<<20 {
		return comm.CodeFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allow := map[string]bool{
		".jpg": true,
		".jpeg": true,
		".png": true,
		".webp": true,
	}
	if !allow[ext] {
		return comm.CodeFileTypeInvalid
	}

	f, err := fileHeader.Open()
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("打开上传文件失败")
		return comm.CodeUploadError
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("读取上传文件失败")
		return comm.CodeUploadError
	}

	sum := md5.Sum(data)
	hash := hex.EncodeToString(sum[:])

	fileRepo := repo.NewFileRecordRepo()
	record, err := fileRepo.FindByHash(ctx, hash)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("查询文件哈希失败")
		return comm.CodeDatabaseError
	}
	if record != nil {
		u.Response.URL = record.FileURL
		u.Response.Hash = hash
		return comm.CodeOK
	}

	ak := qiniuAccessKey
	sk := qiniuSecretKey
	bucket := qiniuBucket
	domain := strings.TrimRight(qiniuDomain, "/")
	if ak == "" || sk == "" || bucket == "" || domain == "" {
		nlog.Pick().WithContext(ctx).Error("qiniu 配置缺失")
		return comm.CodeThirdServiceError
	}

	key := fmt.Sprintf("images/%s%s", hash, ext)
	putPolicy := storage.PutPolicy{Scope: bucket}
	upToken := putPolicy.UploadToken(qbox.NewMac(ak, sk))

	cfg := storage.Config{UseHTTPS: false}
	uploader := storage.NewFormUploader(&cfg)
	var putRet storage.PutRet
	if err = uploader.Put(
		ctx.Request.Context(),
		&putRet,
		upToken,
		key,
		bytes.NewReader(data),
		int64(len(data)),
		nil,
	); err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("七牛上传失败")
		return comm.CodeUploadError
	}

	url := domain + "/" + putRet.Key
	if err = fileRepo.Create(ctx, hash, url); err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("保存文件记录失败")
		return comm.CodeDatabaseError
	}

	u.Response.URL = url
	u.Response.Hash = hash
	return comm.CodeOK
}

// Init Api初始化: 进行参数校验和绑定
func (u *UploadApi) Init(ctx *gin.Context) (err error) {
	if err = ctx.ShouldBindUri(&u.Request.Uri); err != nil {
		return err
	}
	if err = ctx.ShouldBindHeader(&u.Request.Header); err != nil {
		return err
	}
	if err = ctx.ShouldBindQuery(&u.Request.Query); err != nil {
		return err
	}
	if err = ctx.ShouldBind(&u.Request.Form); err != nil {
		return err
	}
	return nil
}

// hfUpload API执行入口
func hfUpload(ctx *gin.Context) {
	api := &UploadApi{}
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

