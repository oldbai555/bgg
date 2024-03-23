// 指定proto版本

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: lbstore.proto

// 指定默认包名

package lbstore

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	lb "github.com/oldbai555/bgg/service/lb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ErrCode int32

const (
	ErrCode_Success ErrCode = 0
	// 100000 - 110000
	ErrCode_ErrFileNotFound ErrCode = 100000 // 文件不存在
)

// Enum value maps for ErrCode.
var (
	ErrCode_name = map[int32]string{
		0:      "Success",
		100000: "ErrFileNotFound",
	}
	ErrCode_value = map[string]int32{
		"Success":         0,
		"ErrFileNotFound": 100000,
	}
)

func (x ErrCode) Enum() *ErrCode {
	p := new(ErrCode)
	*p = x
	return p
}

func (x ErrCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrCode) Descriptor() protoreflect.EnumDescriptor {
	return file_lbstore_proto_enumTypes[0].Descriptor()
}

func (ErrCode) Type() protoreflect.EnumType {
	return &file_lbstore_proto_enumTypes[0]
}

func (x ErrCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrCode.Descriptor instead.
func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{0}
}

type GetFileListReq_Option int32

const (
	GetFileListReq_OptionNil          GetFileListReq_Option = 0
	GetFileListReq_OptionLikeFileName GetFileListReq_Option = 1
)

// Enum value maps for GetFileListReq_Option.
var (
	GetFileListReq_Option_name = map[int32]string{
		0: "OptionNil",
		1: "OptionLikeFileName",
	}
	GetFileListReq_Option_value = map[string]int32{
		"OptionNil":          0,
		"OptionLikeFileName": 1,
	}
)

func (x GetFileListReq_Option) Enum() *GetFileListReq_Option {
	p := new(GetFileListReq_Option)
	*p = x
	return p
}

func (x GetFileListReq_Option) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GetFileListReq_Option) Descriptor() protoreflect.EnumDescriptor {
	return file_lbstore_proto_enumTypes[1].Descriptor()
}

func (GetFileListReq_Option) Type() protoreflect.EnumType {
	return &file_lbstore_proto_enumTypes[1]
}

func (x GetFileListReq_Option) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GetFileListReq_Option.Descriptor instead.
func (GetFileListReq_Option) EnumDescriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{3, 0}
}

type ModelFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @gotags: gorm:"primaryKey"
	Id         uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt  int32  `protobuf:"varint,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt  int32  `protobuf:"varint,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt  int32  `protobuf:"varint,4,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	CreatorUid uint64 `protobuf:"varint,5,opt,name=creator_uid,json=creatorUid,proto3" json:"creator_uid,omitempty"`
	FileName   string `protobuf:"bytes,6,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// @desc: 文件后缀
	FileExt string `protobuf:"bytes,7,opt,name=file_ext,json=fileExt,proto3" json:"file_ext,omitempty"`
	// @desc: 可以用它来续期 signed_url
	ObjectKey string `protobuf:"bytes,8,opt,name=object_key,json=objectKey,proto3" json:"object_key,omitempty"`
	// @desc: 签名的 url
	SignUrl string `protobuf:"bytes,9,opt,name=sign_url,json=signUrl,proto3" json:"sign_url,omitempty"`
	// @desc: 正常的 url
	Url string `protobuf:"bytes,10,opt,name=url,proto3" json:"url,omitempty"`
	// @desc: 文件类型
	FileType string `protobuf:"bytes,11,opt,name=file_type,json=fileType,proto3" json:"file_type,omitempty"`
	// @desc: 文件大小
	Size uint64 `protobuf:"varint,12,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *ModelFile) Reset() {
	*x = ModelFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModelFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModelFile) ProtoMessage() {}

func (x *ModelFile) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModelFile.ProtoReflect.Descriptor instead.
func (*ModelFile) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{0}
}

func (x *ModelFile) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ModelFile) GetCreatedAt() int32 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *ModelFile) GetUpdatedAt() int32 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *ModelFile) GetDeletedAt() int32 {
	if x != nil {
		return x.DeletedAt
	}
	return 0
}

func (x *ModelFile) GetCreatorUid() uint64 {
	if x != nil {
		return x.CreatorUid
	}
	return 0
}

func (x *ModelFile) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *ModelFile) GetFileExt() string {
	if x != nil {
		return x.FileExt
	}
	return ""
}

func (x *ModelFile) GetObjectKey() string {
	if x != nil {
		return x.ObjectKey
	}
	return ""
}

func (x *ModelFile) GetSignUrl() string {
	if x != nil {
		return x.SignUrl
	}
	return ""
}

func (x *ModelFile) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *ModelFile) GetFileType() string {
	if x != nil {
		return x.FileType
	}
	return ""
}

func (x *ModelFile) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

type UploadReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Buf      []byte `protobuf:"bytes,1,opt,name=buf,proto3" json:"buf,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	FileExt  string `protobuf:"bytes,3,opt,name=file_ext,json=fileExt,proto3" json:"file_ext,omitempty"`
}

func (x *UploadReq) Reset() {
	*x = UploadReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadReq) ProtoMessage() {}

func (x *UploadReq) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadReq.ProtoReflect.Descriptor instead.
func (*UploadReq) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{1}
}

func (x *UploadReq) GetBuf() []byte {
	if x != nil {
		return x.Buf
	}
	return nil
}

func (x *UploadReq) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *UploadReq) GetFileExt() string {
	if x != nil {
		return x.FileExt
	}
	return ""
}

type UploadRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *UploadRsp) Reset() {
	*x = UploadRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadRsp) ProtoMessage() {}

func (x *UploadRsp) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadRsp.ProtoReflect.Descriptor instead.
func (*UploadRsp) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{2}
}

func (x *UploadRsp) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type GetFileListReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options *lb.ListOption `protobuf:"bytes,1,opt,name=options,proto3" json:"options,omitempty"`
}

func (x *GetFileListReq) Reset() {
	*x = GetFileListReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileListReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileListReq) ProtoMessage() {}

func (x *GetFileListReq) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileListReq.ProtoReflect.Descriptor instead.
func (*GetFileListReq) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{3}
}

func (x *GetFileListReq) GetOptions() *lb.ListOption {
	if x != nil {
		return x.Options
	}
	return nil
}

type GetFileListRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Paginate *lb.Paginate `protobuf:"bytes,1,opt,name=paginate,proto3" json:"paginate,omitempty"`
	List     []*ModelFile `protobuf:"bytes,2,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *GetFileListRsp) Reset() {
	*x = GetFileListRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileListRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileListRsp) ProtoMessage() {}

func (x *GetFileListRsp) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileListRsp.ProtoReflect.Descriptor instead.
func (*GetFileListRsp) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{4}
}

func (x *GetFileListRsp) GetPaginate() *lb.Paginate {
	if x != nil {
		return x.Paginate
	}
	return nil
}

func (x *GetFileListRsp) GetList() []*ModelFile {
	if x != nil {
		return x.List
	}
	return nil
}

type RefreshFileSignedUrlReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RefreshFileSignedUrlReq) Reset() {
	*x = RefreshFileSignedUrlReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshFileSignedUrlReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshFileSignedUrlReq) ProtoMessage() {}

func (x *RefreshFileSignedUrlReq) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshFileSignedUrlReq.ProtoReflect.Descriptor instead.
func (*RefreshFileSignedUrlReq) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{5}
}

func (x *RefreshFileSignedUrlReq) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type RefreshFileSignedUrlRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RefreshFileSignedUrlRsp) Reset() {
	*x = RefreshFileSignedUrlRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshFileSignedUrlRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshFileSignedUrlRsp) ProtoMessage() {}

func (x *RefreshFileSignedUrlRsp) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshFileSignedUrlRsp.ProtoReflect.Descriptor instead.
func (*RefreshFileSignedUrlRsp) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{6}
}

type GetSignatureReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Method string `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
}

func (x *GetSignatureReq) Reset() {
	*x = GetSignatureReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSignatureReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSignatureReq) ProtoMessage() {}

func (x *GetSignatureReq) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSignatureReq.ProtoReflect.Descriptor instead.
func (*GetSignatureReq) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{7}
}

func (x *GetSignatureReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetSignatureReq) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

type GetSignatureRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signature    string `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	SessionToken string `protobuf:"bytes,2,opt,name=session_token,json=sessionToken,proto3" json:"session_token,omitempty"`
}

func (x *GetSignatureRsp) Reset() {
	*x = GetSignatureRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSignatureRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSignatureRsp) ProtoMessage() {}

func (x *GetSignatureRsp) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSignatureRsp.ProtoReflect.Descriptor instead.
func (*GetSignatureRsp) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{8}
}

func (x *GetSignatureRsp) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *GetSignatureRsp) GetSessionToken() string {
	if x != nil {
		return x.SessionToken
	}
	return ""
}

type ReportUploadFileReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File *ModelFile `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *ReportUploadFileReq) Reset() {
	*x = ReportUploadFileReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportUploadFileReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportUploadFileReq) ProtoMessage() {}

func (x *ReportUploadFileReq) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportUploadFileReq.ProtoReflect.Descriptor instead.
func (*ReportUploadFileReq) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{9}
}

func (x *ReportUploadFileReq) GetFile() *ModelFile {
	if x != nil {
		return x.File
	}
	return nil
}

type ReportUploadFileRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReportUploadFileRsp) Reset() {
	*x = ReportUploadFileRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbstore_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportUploadFileRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportUploadFileRsp) ProtoMessage() {}

func (x *ReportUploadFileRsp) ProtoReflect() protoreflect.Message {
	mi := &file_lbstore_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportUploadFileRsp.ProtoReflect.Descriptor instead.
func (*ReportUploadFileRsp) Descriptor() ([]byte, []int) {
	return file_lbstore_proto_rawDescGZIP(), []int{10}
}

var File_lbstore_proto protoreflect.FileDescriptor

var file_lbstore_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x08, 0x6c, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x0e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xce, 0x02, 0x0a, 0x09, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x46, 0x69, 0x6c, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x55, 0x69, 0x64, 0x12, 0x1b,
	0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x65, 0x78, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x66,
	0x69, 0x6c, 0x65, 0x45, 0x78, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x75, 0x72,
	0x6c, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x69, 0x67, 0x6e, 0x55, 0x72, 0x6c,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x72, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x22, 0x55, 0x0a, 0x09, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71,
	0x12, 0x10, 0x0a, 0x03, 0x62, 0x75, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x62,
	0x75, 0x66, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x45, 0x78, 0x74, 0x22, 0x1d, 0x0a, 0x09, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x73, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x75, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x12, 0x32, 0x0a, 0x07, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6c,
	0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x08, 0xfa, 0x42,
	0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22,
	0x2f, 0x0a, 0x06, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0d, 0x0a, 0x09, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x4e, 0x69, 0x6c, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x4c, 0x69, 0x6b, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x10, 0x01,
	0x22, 0x62, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x73, 0x70, 0x12, 0x28, 0x0a, 0x08, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6c, 0x62, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x65, 0x52, 0x08, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x12, 0x26, 0x0a, 0x04,
	0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6c, 0x62, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04,
	0x6c, 0x69, 0x73, 0x74, 0x22, 0x32, 0x0a, 0x17, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x46,
	0x69, 0x6c, 0x65, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x72, 0x6c, 0x52, 0x65, 0x71, 0x12,
	0x17, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x32, 0x02, 0x20, 0x00, 0x52, 0x02, 0x69, 0x64, 0x22, 0x19, 0x0a, 0x17, 0x52, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x72, 0x6c,
	0x52, 0x73, 0x70, 0x22, 0x4f, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x22, 0x54, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x52, 0x73, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x47, 0x0a, 0x13, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x12, 0x30, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x46,
	0x69, 0x6c, 0x65, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x66,
	0x69, 0x6c, 0x65, 0x22, 0x15, 0x0a, 0x13, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x73, 0x70, 0x2a, 0x2d, 0x0a, 0x07, 0x45, 0x72,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x10, 0x00, 0x12, 0x15, 0x0a, 0x0f, 0x45, 0x72, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x6f, 0x74,
	0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0xa0, 0x8d, 0x06, 0x32, 0x86, 0x03, 0x0a, 0x07, 0x6c, 0x62,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x42, 0x0a, 0x06, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12,
	0x12, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x73, 0x70, 0x22, 0x10, 0x8a, 0xe2, 0x09, 0x04, 0x50, 0x4f, 0x53,
	0x54, 0x92, 0xe2, 0x09, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x41, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x17, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x1a, 0x17, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x46,
	0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x73, 0x70, 0x22, 0x00, 0x12, 0x5c, 0x0a, 0x14,
	0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x55, 0x72, 0x6c, 0x12, 0x20, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x55, 0x72, 0x6c, 0x52, 0x65, 0x71, 0x1a, 0x20, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x55, 0x72, 0x6c, 0x52, 0x73, 0x70, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x18, 0x2e, 0x6c, 0x62, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47,
	0x65, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x73, 0x70, 0x22, 0x00,
	0x12, 0x50, 0x0a, 0x10, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x46, 0x69, 0x6c, 0x65, 0x12, 0x1c, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x1a, 0x1c, 0x2e, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x73, 0x70,
	0x22, 0x00, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6f, 0x6c, 0x64, 0x62, 0x61, 0x69, 0x35, 0x35, 0x35, 0x2f, 0x62, 0x67, 0x67, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x6c, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lbstore_proto_rawDescOnce sync.Once
	file_lbstore_proto_rawDescData = file_lbstore_proto_rawDesc
)

func file_lbstore_proto_rawDescGZIP() []byte {
	file_lbstore_proto_rawDescOnce.Do(func() {
		file_lbstore_proto_rawDescData = protoimpl.X.CompressGZIP(file_lbstore_proto_rawDescData)
	})
	return file_lbstore_proto_rawDescData
}

var file_lbstore_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_lbstore_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_lbstore_proto_goTypes = []interface{}{
	(ErrCode)(0),                    // 0: lbstore.ErrCode
	(GetFileListReq_Option)(0),      // 1: lbstore.GetFileListReq.Option
	(*ModelFile)(nil),               // 2: lbstore.ModelFile
	(*UploadReq)(nil),               // 3: lbstore.UploadReq
	(*UploadRsp)(nil),               // 4: lbstore.UploadRsp
	(*GetFileListReq)(nil),          // 5: lbstore.GetFileListReq
	(*GetFileListRsp)(nil),          // 6: lbstore.GetFileListRsp
	(*RefreshFileSignedUrlReq)(nil), // 7: lbstore.RefreshFileSignedUrlReq
	(*RefreshFileSignedUrlRsp)(nil), // 8: lbstore.RefreshFileSignedUrlRsp
	(*GetSignatureReq)(nil),         // 9: lbstore.GetSignatureReq
	(*GetSignatureRsp)(nil),         // 10: lbstore.GetSignatureRsp
	(*ReportUploadFileReq)(nil),     // 11: lbstore.ReportUploadFileReq
	(*ReportUploadFileRsp)(nil),     // 12: lbstore.ReportUploadFileRsp
	(*lb.ListOption)(nil),           // 13: lb.ListOption
	(*lb.Paginate)(nil),             // 14: lb.Paginate
}
var file_lbstore_proto_depIdxs = []int32{
	13, // 0: lbstore.GetFileListReq.options:type_name -> lb.ListOption
	14, // 1: lbstore.GetFileListRsp.paginate:type_name -> lb.Paginate
	2,  // 2: lbstore.GetFileListRsp.list:type_name -> lbstore.ModelFile
	2,  // 3: lbstore.ReportUploadFileReq.file:type_name -> lbstore.ModelFile
	3,  // 4: lbstore.lbstore.Upload:input_type -> lbstore.UploadReq
	5,  // 5: lbstore.lbstore.GetFileList:input_type -> lbstore.GetFileListReq
	7,  // 6: lbstore.lbstore.RefreshFileSignedUrl:input_type -> lbstore.RefreshFileSignedUrlReq
	9,  // 7: lbstore.lbstore.GetSignature:input_type -> lbstore.GetSignatureReq
	11, // 8: lbstore.lbstore.ReportUploadFile:input_type -> lbstore.ReportUploadFileReq
	4,  // 9: lbstore.lbstore.Upload:output_type -> lbstore.UploadRsp
	6,  // 10: lbstore.lbstore.GetFileList:output_type -> lbstore.GetFileListRsp
	8,  // 11: lbstore.lbstore.RefreshFileSignedUrl:output_type -> lbstore.RefreshFileSignedUrlRsp
	10, // 12: lbstore.lbstore.GetSignature:output_type -> lbstore.GetSignatureRsp
	12, // 13: lbstore.lbstore.ReportUploadFile:output_type -> lbstore.ReportUploadFileRsp
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_lbstore_proto_init() }
func file_lbstore_proto_init() {
	if File_lbstore_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lbstore_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ModelFile); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileListReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileListRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshFileSignedUrlReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshFileSignedUrlRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSignatureReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSignatureRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportUploadFileReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbstore_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportUploadFileRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_lbstore_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_lbstore_proto_goTypes,
		DependencyIndexes: file_lbstore_proto_depIdxs,
		EnumInfos:         file_lbstore_proto_enumTypes,
		MessageInfos:      file_lbstore_proto_msgTypes,
	}.Build()
	File_lbstore_proto = out.File
	file_lbstore_proto_rawDesc = nil
	file_lbstore_proto_goTypes = nil
	file_lbstore_proto_depIdxs = nil
}
