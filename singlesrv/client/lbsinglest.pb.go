// 指定proto版本

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: lbsinglest.proto

// 指定默认包名

package client

import (
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

// 消息内容
type Content struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text     *Content_Text     `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Image    *Content_Image    `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
	Video    *Content_Video    `protobuf:"bytes,3,opt,name=video,proto3" json:"video,omitempty"`
	Voice    *Content_Voice    `protobuf:"bytes,4,opt,name=voice,proto3" json:"voice,omitempty"`
	Document *Content_Document `protobuf:"bytes,5,opt,name=document,proto3" json:"document,omitempty"`
	Location *Content_Location `protobuf:"bytes,6,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *Content) Reset() {
	*x = Content{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbsinglest_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content) ProtoMessage() {}

func (x *Content) ProtoReflect() protoreflect.Message {
	mi := &file_lbsinglest_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content.ProtoReflect.Descriptor instead.
func (*Content) Descriptor() ([]byte, []int) {
	return file_lbsinglest_proto_rawDescGZIP(), []int{0}
}

func (x *Content) GetText() *Content_Text {
	if x != nil {
		return x.Text
	}
	return nil
}

func (x *Content) GetImage() *Content_Image {
	if x != nil {
		return x.Image
	}
	return nil
}

func (x *Content) GetVideo() *Content_Video {
	if x != nil {
		return x.Video
	}
	return nil
}

func (x *Content) GetVoice() *Content_Voice {
	if x != nil {
		return x.Voice
	}
	return nil
}

func (x *Content) GetDocument() *Content_Document {
	if x != nil {
		return x.Document
	}
	return nil
}

func (x *Content) GetLocation() *Content_Location {
	if x != nil {
		return x.Location
	}
	return nil
}

// 基础用户信息
type BaseUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Avatar   string `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Nickname string `protobuf:"bytes,4,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Email    string `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
	Github   string `protobuf:"bytes,6,opt,name=github,proto3" json:"github,omitempty"`
	Desc     string `protobuf:"bytes,7,opt,name=desc,proto3" json:"desc,omitempty"`
	Role     uint32 `protobuf:"varint,8,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *BaseUser) Reset() {
	*x = BaseUser{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbsinglest_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseUser) ProtoMessage() {}

func (x *BaseUser) ProtoReflect() protoreflect.Message {
	mi := &file_lbsinglest_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseUser.ProtoReflect.Descriptor instead.
func (*BaseUser) Descriptor() ([]byte, []int) {
	return file_lbsinglest_proto_rawDescGZIP(), []int{1}
}

func (x *BaseUser) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *BaseUser) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *BaseUser) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *BaseUser) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *BaseUser) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *BaseUser) GetGithub() string {
	if x != nil {
		return x.Github
	}
	return ""
}

func (x *BaseUser) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *BaseUser) GetRole() uint32 {
	if x != nil {
		return x.Role
	}
	return 0
}

type Content_Text struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	// 翻译结果
	Translated string `protobuf:"bytes,2,opt,name=translated,proto3" json:"translated,omitempty"`
}

func (x *Content_Text) Reset() {
	*x = Content_Text{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbsinglest_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content_Text) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content_Text) ProtoMessage() {}

func (x *Content_Text) ProtoReflect() protoreflect.Message {
	mi := &file_lbsinglest_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content_Text.ProtoReflect.Descriptor instead.
func (*Content_Text) Descriptor() ([]byte, []int) {
	return file_lbsinglest_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Content_Text) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Content_Text) GetTranslated() string {
	if x != nil {
		return x.Translated
	}
	return ""
}

type Content_Image struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// 描述
	Caption string `protobuf:"bytes,2,opt,name=caption,proto3" json:"caption,omitempty"`
	// 名称
	FileName string `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// 格式
	Format string `protobuf:"bytes,4,opt,name=format,proto3" json:"format,omitempty"`
}

func (x *Content_Image) Reset() {
	*x = Content_Image{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbsinglest_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content_Image) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content_Image) ProtoMessage() {}

func (x *Content_Image) ProtoReflect() protoreflect.Message {
	mi := &file_lbsinglest_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content_Image.ProtoReflect.Descriptor instead.
func (*Content_Image) Descriptor() ([]byte, []int) {
	return file_lbsinglest_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Content_Image) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Content_Image) GetCaption() string {
	if x != nil {
		return x.Caption
	}
	return ""
}

func (x *Content_Image) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *Content_Image) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

type Content_Video struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// 描述
	Caption string `protobuf:"bytes,2,opt,name=caption,proto3" json:"caption,omitempty"`
	// 名称
	FileName string `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// 格式
	Format string `protobuf:"bytes,4,opt,name=format,proto3" json:"format,omitempty"`
	// 视频封面
	CoverUrl string `protobuf:"bytes,5,opt,name=cover_url,json=coverUrl,proto3" json:"cover_url,omitempty"`
}

func (x *Content_Video) Reset() {
	*x = Content_Video{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbsinglest_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content_Video) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content_Video) ProtoMessage() {}

func (x *Content_Video) ProtoReflect() protoreflect.Message {
	mi := &file_lbsinglest_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content_Video.ProtoReflect.Descriptor instead.
func (*Content_Video) Descriptor() ([]byte, []int) {
	return file_lbsinglest_proto_rawDescGZIP(), []int{0, 2}
}

func (x *Content_Video) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Content_Video) GetCaption() string {
	if x != nil {
		return x.Caption
	}
	return ""
}

func (x *Content_Video) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *Content_Video) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *Content_Video) GetCoverUrl() string {
	if x != nil {
		return x.CoverUrl
	}
	return ""
}

type Content_Voice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// 描述
	Caption string `protobuf:"bytes,2,opt,name=caption,proto3" json:"caption,omitempty"`
	// 名称
	FileName string `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// 格式
	Format string `protobuf:"bytes,4,opt,name=format,proto3" json:"format,omitempty"`
	// 语音识别结果
	Recognition string `protobuf:"bytes,5,opt,name=recognition,proto3" json:"recognition,omitempty"`
}

func (x *Content_Voice) Reset() {
	*x = Content_Voice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbsinglest_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content_Voice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content_Voice) ProtoMessage() {}

func (x *Content_Voice) ProtoReflect() protoreflect.Message {
	mi := &file_lbsinglest_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content_Voice.ProtoReflect.Descriptor instead.
func (*Content_Voice) Descriptor() ([]byte, []int) {
	return file_lbsinglest_proto_rawDescGZIP(), []int{0, 3}
}

func (x *Content_Voice) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Content_Voice) GetCaption() string {
	if x != nil {
		return x.Caption
	}
	return ""
}

func (x *Content_Voice) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *Content_Voice) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *Content_Voice) GetRecognition() string {
	if x != nil {
		return x.Recognition
	}
	return ""
}

type Content_Document struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// 描述
	Caption string `protobuf:"bytes,2,opt,name=caption,proto3" json:"caption,omitempty"`
	// 名称
	FileName string `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// 格式
	Format string `protobuf:"bytes,4,opt,name=format,proto3" json:"format,omitempty"`
}

func (x *Content_Document) Reset() {
	*x = Content_Document{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbsinglest_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content_Document) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content_Document) ProtoMessage() {}

func (x *Content_Document) ProtoReflect() protoreflect.Message {
	mi := &file_lbsinglest_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content_Document.ProtoReflect.Descriptor instead.
func (*Content_Document) Descriptor() ([]byte, []int) {
	return file_lbsinglest_proto_rawDescGZIP(), []int{0, 4}
}

func (x *Content_Document) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Content_Document) GetCaption() string {
	if x != nil {
		return x.Caption
	}
	return ""
}

func (x *Content_Document) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *Content_Document) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

type Content_Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 地理位置纬度
	X float64 `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	// 地理位置经度
	Y float64 `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
	// 地图缩放大小
	Scale float64 `protobuf:"fixed64,3,opt,name=scale,proto3" json:"scale,omitempty"`
	// 地理位置信息
	Label string `protobuf:"bytes,4,opt,name=label,proto3" json:"label,omitempty"`
}

func (x *Content_Location) Reset() {
	*x = Content_Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbsinglest_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content_Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content_Location) ProtoMessage() {}

func (x *Content_Location) ProtoReflect() protoreflect.Message {
	mi := &file_lbsinglest_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content_Location.ProtoReflect.Descriptor instead.
func (*Content_Location) Descriptor() ([]byte, []int) {
	return file_lbsinglest_proto_rawDescGZIP(), []int{0, 5}
}

func (x *Content_Location) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Content_Location) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *Content_Location) GetScale() float64 {
	if x != nil {
		return x.Scale
	}
	return 0
}

func (x *Content_Location) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

var File_lbsinglest_proto protoreflect.FileDescriptor

var file_lbsinglest_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6c, 0x62, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x6c, 0x62, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x73, 0x74, 0x22, 0xc0,
	0x07, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x2c, 0x0a, 0x04, 0x74, 0x65,
	0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6c, 0x62, 0x73, 0x69, 0x6e,
	0x67, 0x6c, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x54, 0x65,
	0x78, 0x74, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x2f, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6c, 0x62, 0x73, 0x69, 0x6e, 0x67,
	0x6c, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x2f, 0x0a, 0x05, 0x76, 0x69, 0x64,
	0x65, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6c, 0x62, 0x73, 0x69, 0x6e,
	0x67, 0x6c, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x56, 0x69,
	0x64, 0x65, 0x6f, 0x52, 0x05, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x12, 0x2f, 0x0a, 0x05, 0x76, 0x6f,
	0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6c, 0x62, 0x73, 0x69,
	0x6e, 0x67, 0x6c, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x56,
	0x6f, 0x69, 0x63, 0x65, 0x52, 0x05, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x08, 0x64,
	0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x6c, 0x62, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x64, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6c, 0x62, 0x73, 0x69, 0x6e, 0x67,
	0x6c, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a,
	0x40, 0x0a, 0x04, 0x54, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65,
	0x64, 0x1a, 0x68, 0x0a, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x61, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63,
	0x61, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x1a, 0x85, 0x01, 0x0a, 0x05,
	0x56, 0x69, 0x64, 0x65, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x61, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x61, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x55, 0x72, 0x6c, 0x1a, 0x8a, 0x01, 0x0a, 0x05, 0x56, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x61, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x61, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x20,
	0x0a, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x67, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x67, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x1a, 0x6b, 0x0a, 0x08, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x61, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x61, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x1a, 0x52, 0x0a,
	0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x01, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x22, 0xc0, 0x01, 0x0a, 0x08, 0x42, 0x61, 0x73, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x76,
	0x61, 0x74, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x76, 0x61, 0x74,
	0x61, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x65, 0x73, 0x63, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63,
	0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04,
	0x72, 0x6f, 0x6c, 0x65, 0x32, 0x0c, 0x0a, 0x0a, 0x6c, 0x62, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65,
	0x73, 0x74, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6f, 0x6c, 0x64, 0x62, 0x61, 0x69, 0x35, 0x35, 0x35, 0x2f, 0x62, 0x67, 0x67, 0x2f, 0x73,
	0x69, 0x6e, 0x67, 0x6c, 0x65, 0x73, 0x72, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lbsinglest_proto_rawDescOnce sync.Once
	file_lbsinglest_proto_rawDescData = file_lbsinglest_proto_rawDesc
)

func file_lbsinglest_proto_rawDescGZIP() []byte {
	file_lbsinglest_proto_rawDescOnce.Do(func() {
		file_lbsinglest_proto_rawDescData = protoimpl.X.CompressGZIP(file_lbsinglest_proto_rawDescData)
	})
	return file_lbsinglest_proto_rawDescData
}

var file_lbsinglest_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_lbsinglest_proto_goTypes = []interface{}{
	(*Content)(nil),          // 0: lbsinglest.Content
	(*BaseUser)(nil),         // 1: lbsinglest.BaseUser
	(*Content_Text)(nil),     // 2: lbsinglest.Content.Text
	(*Content_Image)(nil),    // 3: lbsinglest.Content.Image
	(*Content_Video)(nil),    // 4: lbsinglest.Content.Video
	(*Content_Voice)(nil),    // 5: lbsinglest.Content.Voice
	(*Content_Document)(nil), // 6: lbsinglest.Content.Document
	(*Content_Location)(nil), // 7: lbsinglest.Content.Location
}
var file_lbsinglest_proto_depIdxs = []int32{
	2, // 0: lbsinglest.Content.text:type_name -> lbsinglest.Content.Text
	3, // 1: lbsinglest.Content.image:type_name -> lbsinglest.Content.Image
	4, // 2: lbsinglest.Content.video:type_name -> lbsinglest.Content.Video
	5, // 3: lbsinglest.Content.voice:type_name -> lbsinglest.Content.Voice
	6, // 4: lbsinglest.Content.document:type_name -> lbsinglest.Content.Document
	7, // 5: lbsinglest.Content.location:type_name -> lbsinglest.Content.Location
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_lbsinglest_proto_init() }
func file_lbsinglest_proto_init() {
	if File_lbsinglest_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lbsinglest_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content); i {
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
		file_lbsinglest_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseUser); i {
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
		file_lbsinglest_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content_Text); i {
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
		file_lbsinglest_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content_Image); i {
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
		file_lbsinglest_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content_Video); i {
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
		file_lbsinglest_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content_Voice); i {
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
		file_lbsinglest_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content_Document); i {
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
		file_lbsinglest_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content_Location); i {
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
			RawDescriptor: file_lbsinglest_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_lbsinglest_proto_goTypes,
		DependencyIndexes: file_lbsinglest_proto_depIdxs,
		MessageInfos:      file_lbsinglest_proto_msgTypes,
	}.Build()
	File_lbsinglest_proto = out.File
	file_lbsinglest_proto_rawDesc = nil
	file_lbsinglest_proto_goTypes = nil
	file_lbsinglest_proto_depIdxs = nil
}