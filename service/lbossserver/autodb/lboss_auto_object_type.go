package autodb

import (
	"github.com/oldbai555/micro/gormx/engine"
)

var ModelFile = &engine.ModelObjectType{
	Name: "lboss.ModelFile",
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
				FieldName: "bucket",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "domain",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}
