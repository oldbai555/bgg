message 生成示例
```go
type ModelUser struct {
    Id                   uint64   `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
    CreatedAt            int32    `protobuf:"varint,2,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
    UpdatedAt            int32    `protobuf:"varint,3,opt,name=updated_at,json=updatedAt" json:"updated_at,omitempty"`
    DeletedAt            int32    `protobuf:"varint,4,opt,name=deleted_at,json=deletedAt" json:"deleted_at,omitempty"`
    Username             string   `protobuf:"bytes,5,opt,name=username" json:"username,omitempty"`
    Password             string   `protobuf:"bytes,6,opt,name=password" json:"password,omitempty"`
    Avatar               string   `protobuf:"bytes,7,opt,name=avatar" json:"avatar,omitempty"`
    Nickname             string   `protobuf:"bytes,8,opt,name=nickname" json:"nickname,omitempty"`
    Email                string   `protobuf:"bytes,9,opt,name=email" json:"email,omitempty"`
    Github               string   `protobuf:"bytes,10,opt,name=github" json:"github,omitempty"`
    Desc                 string   `protobuf:"bytes,11,opt,name=desc" json:"desc,omitempty"`
    Role                 uint32   `protobuf:"varint,12,opt,name=role" json:"role,omitempty"`
    XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-"`
    XXX_unrecognized     []byte   `json:"-" bson:"-"`
    XXX_sizecache        int32    `json:"-" bson:"-"`
}

```