package autodb

import (
	"github.com/oldbai555/micro/gormx/engine"
)

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

var ModelChat = &engine.ModelObjectType{
	Name: "lbsingle.ModelChat",
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
				FieldName: "type",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "account_sn",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "customer_sn",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMessage = &engine.ModelObjectType{
	Name: "lbsingle.ModelMessage",
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
				FieldName: "server_msg_id",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "sys_msg_id",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "send_at",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "from",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "to",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "content",
				Type:      "object",
				IsArray:   false,
			},
			{
				FieldName: "source",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "type",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "status",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "chat_id",
				Type:      "uint64",
				IsArray:   false,
			},
		},
	},
}

var ModelDailyShortSentences = &engine.ModelObjectType{
	Name: "lbsingle.ModelDailyShortSentences",
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
				FieldName: "type",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "content",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "img",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "literature_author",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "convert_img",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelOutsideWebSite = &engine.ModelObjectType{
	Name: "lbsingle.ModelOutsideWebSite",
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
				FieldName: "url",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "title",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "desc",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}
