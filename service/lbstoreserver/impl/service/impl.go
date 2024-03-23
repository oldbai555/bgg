package service

import (
	"context"
	"github.com/oldbai555/bgg/service/lbstore"
)

var OnceSvrImpl = &LbstoreServer{}

type LbstoreServer struct {
	lbstore.UnimplementedLbstoreServer
}

func (a *LbstoreServer) Upload(ctx context.Context, req *lbstore.UploadReq) (*lbstore.UploadRsp, error) {
	var rsp lbstore.UploadRsp
	var err error

	return &rsp, err
}
func (a *LbstoreServer) GetFileList(ctx context.Context, req *lbstore.GetFileListReq) (*lbstore.GetFileListRsp, error) {
	var rsp lbstore.GetFileListRsp
	var err error

	return &rsp, err
}
func (a *LbstoreServer) RefreshFileSignedUrl(ctx context.Context, req *lbstore.RefreshFileSignedUrlReq) (*lbstore.RefreshFileSignedUrlRsp, error) {
	var rsp lbstore.RefreshFileSignedUrlRsp
	var err error

	return &rsp, err
}
func (a *LbstoreServer) GetSignature(ctx context.Context, req *lbstore.GetSignatureReq) (*lbstore.GetSignatureRsp, error) {
	var rsp lbstore.GetSignatureRsp
	var err error

	return &rsp, err
}
func (a *LbstoreServer) ReportUploadFile(ctx context.Context, req *lbstore.ReportUploadFileReq) (*lbstore.ReportUploadFileRsp, error) {
	var rsp lbstore.ReportUploadFileRsp
	var err error

	return &rsp, err
}
