import type * as lbstore from "./model/lbstore";
import api from "./config/api";

async function upload(req: lbstore.UploadReq): Promise<lbstore.UploadRsp> {
    const response = await api.post<lbstore.UploadRsp>("/lbstore/Upload", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getFileList(req: lbstore.GetFileListReq): Promise<lbstore.GetFileListRsp> {
    const response = await api.post<lbstore.GetFileListRsp>("/lbstore/GetFileList", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function refreshFileSignedUrl(req: lbstore.RefreshFileSignedUrlReq): Promise<lbstore.RefreshFileSignedUrlRsp> {
    const response = await api.post<lbstore.RefreshFileSignedUrlRsp>("/lbstore/RefreshFileSignedUrl", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getSignature(req: lbstore.GetSignatureReq): Promise<lbstore.GetSignatureRsp> {
    const response = await api.post<lbstore.GetSignatureRsp>("/lbstore/GetSignature", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function reportUploadFile(req: lbstore.ReportUploadFileReq): Promise<lbstore.ReportUploadFileRsp> {
    const response = await api.post<lbstore.ReportUploadFileRsp>("/lbstore/ReportUploadFile", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

export default {
	upload,
	getFileList,
	refreshFileSignedUrl,
	getSignature,
	reportUploadFile
}
