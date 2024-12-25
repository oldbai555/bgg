package autodb

import (
	"github.com/oldbai555/micro/gormx/engine"
)

var ModelMpMemberUser = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpMemberUser",
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
				FieldName: "nickname",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "avatar",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "status",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "mobile",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "register_ip",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "last_login_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "real_name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "birthday",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "card_id",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "mark",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "last_login_ip",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "now_money",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "sign_num",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "level",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "integral",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "pay_count",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "login_type",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "mp_openid",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "gender",
				Type:      "uint32",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreShop = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpStoreShop",
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
				FieldName: "name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "mobile",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "image",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "images",
				Type:      "string",
				IsArray:   true,
			},
			{
				FieldName: "address",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "notice",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "status",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "lng",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "lat",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpProductCategory = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpProductCategory",
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
				FieldName: "mp_store_shop_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "pic_url",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "description",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreProduct = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpStoreProduct",
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
				FieldName: "mp_store_shop_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "cate_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "image",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "slider_image",
				Type:      "string",
				IsArray:   true,
			},
			{
				FieldName: "name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "info",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreOrderCartInfo = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpStoreOrderCartInfo",
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
				FieldName: "mp_order_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "order_sn",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "product_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "title",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "image",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "number",
				Type:      "uint32",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreOrder = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpStoreOrder",
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
				FieldName: "order_sn",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "mp_store_shop_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "number_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "total_num",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "pay_type",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "remark",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpService = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpService",
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
				FieldName: "name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "image",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "type",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "content",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "pid",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "app_id",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "pages",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "phone",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "weigh",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "status",
				Type:      "int32",
				IsArray:   false,
			},
		},
	},
}

var ModelMpShopAds = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpShopAds",
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
				FieldName: "image",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "is_switch",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "weigh",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "mp_store_shop_id",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreOrderStatus = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpStoreOrderStatus",
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
				FieldName: "oid",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "change_type",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "change_message",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpOrderNumber = &engine.ModelObjectType{
	Name: "lbwxmp.ModelMpOrderNumber",
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
				FieldName: "order_sn",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}
