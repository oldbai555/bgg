package service

import (
	"bytes"
	"context"
	"github.com/oldbai555/bgg/pkg/storage"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/bgg/service/lbstore"
	"github.com/oldbai555/bgg/service/lbstoreserver/impl/dao/impl/mysql"
	"github.com/oldbai555/lbtool/log"
	"net/http"
	"time"
)

var ServerImpl LbstoreServer

type LbstoreServer struct {
	*lbstore.UnimplementedLbstoreServer
}

const defaultExpiredInSec = 60 * 60 * 24 * 365
const defaultObjectKeyPrefix = `public/link/assets/file/`

func (a *LbstoreServer) Upload(ctx context.Context, req *lbstore.UploadReq) (*lbstore.UploadRsp, error) {
	var rsp lbstore.UploadRsp
	var err error

	var objectKey = defaultObjectKeyPrefix + req.FileName

	err = storage.S.Put(objectKey, bytes.NewReader(req.Buf))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	signURL, err := storage.S.SignURL(objectKey, http.MethodGet, defaultExpiredInSec)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.Url = signURL

	_, err = mysql.FileOrm.Create(ctx, &lbstore.ModelFile{
		FileName:  req.FileName,
		FileExt:   req.FileExt,
		ObjectKey: objectKey,
		SignUrl:   signURL,
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbstoreServer) GetFileList(ctx context.Context, req *lbstore.GetFileListReq) (*lbstore.GetFileListRsp, error) {
	var rsp lbstore.GetFileListRsp
	var err error

	rsp.List, rsp.Paginate, err = mysql.FileOrm.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbstoreServer) RefreshFileSignedUrl(ctx context.Context, req *lbstore.RefreshFileSignedUrlReq) (*lbstore.RefreshFileSignedUrlRsp, error) {
	var rsp lbstore.RefreshFileSignedUrlRsp
	var err error

	file, err := mysql.FileOrm.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	signURL, err := storage.S.SignURL(file.ObjectKey, http.MethodGet, defaultExpiredInSec)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	_, err = mysql.FileOrm.UpdateById(ctx, file.Id, map[string]interface{}{
		lbstore.FieldSignUrl_: signURL,
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbstoreServer) GetSignature(ctx context.Context, req *lbstore.GetSignatureReq) (*lbstore.GetSignatureRsp, error) {
	var rsp lbstore.GetSignatureRsp
	var err error

	if req.Method == "" {
		req.Method = http.MethodPut
	}

	credentials, err := storage.S.GetCredentials()
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.SessionToken = credentials.SessionToken

	signature := storage.S.GetSignature(req.Method, req.Name, credentials.SecretID, credentials.SecretKey, time.Minute)
	rsp.Signature = signature

	return &rsp, err
}
func (a *LbstoreServer) ReportUploadFile(ctx context.Context, req *lbstore.ReportUploadFileReq) (*lbstore.ReportUploadFileRsp, error) {
	var rsp lbstore.ReportUploadFileRsp
	var err error

	claims, err := webtool.GetClaimsWithCtx(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	req.File.CreatorUid = claims.GetUserId()
	_, err = mysql.FileOrm.Create(ctx, req.File)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, err
}
