
// Code generated by gen_sys_cmd_autogen.go, DO NOT EDIT.

import type * as lb from "./lb"

export interface ModelFile {
	id : number | undefined;
	created_at : number | undefined;
	updated_at : number | undefined;
	deleted_at : number | undefined;
	creator_uid : number | undefined;
	file_name : string | undefined;
	file_ext : string | undefined;
	object_key : string | undefined;
	sign_url : string | undefined;
	url : string | undefined;
	file_type : string | undefined;
	size : number | undefined;
}


export interface ModelFile {
	id : number | undefined;
	created_at : number | undefined;
	updated_at : number | undefined;
	deleted_at : number | undefined;
	creator_uid : number | undefined;
	file_name : string | undefined;
	file_ext : string | undefined;
	object_key : string | undefined;
	sign_url : string | undefined;
	url : string | undefined;
	file_type : string | undefined;
	size : number | undefined;
}


export interface UploadReq {
	buf : Uint8Array | undefined;
	file_name : string | undefined;
	file_ext : string | undefined;
}


export interface UploadReq {
	buf : Uint8Array | undefined;
	file_name : string | undefined;
	file_ext : string | undefined;
}


export interface UploadRsp {
	url : string | undefined;
}


export interface UploadRsp {
	url : string | undefined;
}


export interface GetFileListReq {
	options : lb.Options | undefined;
}


export interface GetFileListReq {
	options : lb.Options | undefined;
}


export interface GetFileListRsp {
	paginate : lb.Paginate | undefined;
	list : ModelFile[];
}


export interface GetFileListRsp {
	paginate : lb.Paginate | undefined;
	list : ModelFile[];
}


export interface RefreshFileSignedUrlReq {
	id : number | undefined;
}


export interface RefreshFileSignedUrlReq {
	id : number | undefined;
}


export interface RefreshFileSignedUrlRsp {
}


export interface RefreshFileSignedUrlRsp {
}


export interface GetSignatureReq {
	name : string | undefined;
	method : string | undefined;
}


export interface GetSignatureReq {
	name : string | undefined;
	method : string | undefined;
}


export interface GetSignatureRsp {
	signature : string | undefined;
	session_token : string | undefined;
}


export interface GetSignatureRsp {
	signature : string | undefined;
	session_token : string | undefined;
}


export interface ReportUploadFileReq {
	file : ModelFile | undefined;
}


export interface ReportUploadFileReq {
	file : ModelFile | undefined;
}


export interface ReportUploadFileRsp {
}


export interface ReportUploadFileRsp {
}


export enum ErrCode {
	Success=0,
}

export enum GetFileListReq_Option {
	OptionNil=0,
	OptionLikeFileName=1,
}

