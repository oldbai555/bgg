package lbsingledb

import (
	"github.com/oldbai555/micro/gormx/engine"
)

var ModelFile = &engine.ModelObjectType{
	Name: "lbsingle.ModelFile",
	FieldList: &engine.ObjectFieldList{
		List: []*engine.ObjectField{
			{
				FieldName: "id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "created_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "updated_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "deleted_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "creator_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "size",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "rename",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "path",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "md5",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "sort_url",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "state",
				Type:      "uint32",
				IsArray:   false,
			},
		},
	},
}

var ModelUser = &engine.ModelObjectType{
	Name: "lbsingle.ModelUser",
	FieldList: &engine.ObjectFieldList{
		List: []*engine.ObjectField{
			{
				FieldName: "id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "created_at",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "updated_at",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "deleted_at",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "username",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "password",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "avatar",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "nickname",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "email",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "github",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "desc",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "role",
				Type:      "uint32",
				IsArray:   false,
			},
		},
	},
}
