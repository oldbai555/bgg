package lbossserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/pkg/bctx"
	"github.com/oldbai555/bgg/pkg/constant"
	"github.com/oldbai555/bgg/service/lboss"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bgin"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/uctx"
	"path"
	"strings"
	"time"
)

func (a *LbossServer) DelFileList(ctx context.Context, req *lboss.DelFileListReq) (*lboss.DelFileListRsp, error) {
	var rsp lboss.DelFileListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetFileList(ctx, &lboss.GetFileListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, lboss.FieldId_),
	})
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lboss.FieldId)
	_, err = OrmFile.NewBaseScope().WhereIn(lboss.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbossServer) GetFile(ctx context.Context, req *lboss.GetFileReq) (*lboss.GetFileRsp, error) {
	var rsp lboss.GetFileRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmFile.NewBaseScope().Where(lboss.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data

	return &rsp, err
}

// 预签名
func handlePreSigned(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)
	if minIoSDK == nil {
		return
	}
	nCtx := bctx.NewCtx(ctx, bctx.WithGinHeaderAuthorization(ctx))
	_, err := lbsingle.CheckAuth(nCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	object, err := minIoSDK.PreSignedPutObject(constant.BucketByPublic, ctx.Param("fileName"))
	if err != nil {
		handler.Error(err)
		return
	}
	handler.HttpJson(object)
}

// 上传文件
func handleUploadFile(c *gin.Context) {
	handler := bgin.NewHandler(c)
	nCtx := bctx.NewCtx(c, bctx.WithGinHeaderAuthorization(c))
	_, err := lbsingle.CheckAuth(nCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	open, err := file.Open()
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	md5 := utils.StrMd5(file.Filename)
	reFileName := md5 + path.Ext(file.Filename)
	timeFmt := time.Unix(time.Now().Unix(), 0).Format("20060102")
	oss, err := uploadToOss(nCtx, file.Size, timeFmt, file.Filename, reFileName, open)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	handler.HttpJson(oss)
}

// 上传包
func handleUploadDeployFile(c *gin.Context) {
	handler := bgin.NewHandler(c)
	nCtx := bctx.NewCtx(c)
	// 获取文件信息
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	open, err := file.Open()
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	oss, err := uploadToOss(nCtx, file.Size, "deploy", file.Filename, file.Filename, open)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	handler.HttpJson(oss)
}

// 下载文件
func handleDownloadFile(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)
	p := ctx.Param("path")
	nCtx := bctx.NewCtx(ctx)
	modelFile, err := OrmFile.NewBaseScope().Select(lboss.FieldPath_).Where(lboss.FieldPath_, strings.TrimLeft(p, "/")).First(nCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	u, err := minIoSDK.PreSignedGetObject(constant.BucketByPublic, modelFile.Path)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(lboss.ErrFileNotFound)
		return
	}
	ctx.Redirect(301, u)
}
