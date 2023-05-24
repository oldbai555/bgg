message 生成示例
```go
type ModelUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt int32  `protobuf:"varint,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt int32  `protobuf:"varint,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt int32  `protobuf:"varint,4,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	Username  string `protobuf:"bytes,5,opt,name=username,proto3" json:"username,omitempty"`
	Password  string `protobuf:"bytes,6,opt,name=password,proto3" json:"password,omitempty"`
	Avatar    string `protobuf:"bytes,7,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Nickname  string `protobuf:"bytes,8,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Email     string `protobuf:"bytes,9,opt,name=email,proto3" json:"email,omitempty"`
	Github    string `protobuf:"bytes,10,opt,name=github,proto3" json:"github,omitempty"`
	Desc      string `protobuf:"bytes,11,opt,name=desc,proto3" json:"desc,omitempty"`
	Role      uint32 `protobuf:"varint,12,opt,name=role,proto3" json:"role,omitempty"`
}
```