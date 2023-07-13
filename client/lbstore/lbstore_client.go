// Code generated by gen_client.go, DO NOT EDIT.
package lbstore

import (
	"context"
)

func Upload(ctx context.Context, req *UploadReq) (*UploadRsp, error) {
	if cliMgr.conn == nil {
		return nil, cliMgr.Err
	}
	return cliMgr.client.Upload(ctx, req)
}

func GetFileList(ctx context.Context, req *GetFileListReq) (*GetFileListRsp, error) {
	if cliMgr.conn == nil {
		return nil, cliMgr.Err
	}
	return cliMgr.client.GetFileList(ctx, req)
}

func RefreshFileSignedUrl(ctx context.Context, req *RefreshFileSignedUrlReq) (*RefreshFileSignedUrlRsp, error) {
	if cliMgr.conn == nil {
		return nil, cliMgr.Err
	}
	return cliMgr.client.RefreshFileSignedUrl(ctx, req)
}

func GetSignature(ctx context.Context, req *GetSignatureReq) (*GetSignatureRsp, error) {
	if cliMgr.conn == nil {
		return nil, cliMgr.Err
	}
	return cliMgr.client.GetSignature(ctx, req)
}

func ReportUploadFile(ctx context.Context, req *ReportUploadFileReq) (*ReportUploadFileRsp, error) {
	if cliMgr.conn == nil {
		return nil, cliMgr.Err
	}
	return cliMgr.client.ReportUploadFile(ctx, req)
}
