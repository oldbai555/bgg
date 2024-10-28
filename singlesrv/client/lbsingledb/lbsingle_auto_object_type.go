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

var ModelMpMerchantDetails = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpMerchantDetails",
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
		},
	},
}

var ModelMpMemberUser = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpMemberUser",
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

var ModelMpUserAddress = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpUserAddress",
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
				FieldName: "mp_uid",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "real_name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "phone",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "province",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "city",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "cityId",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "district",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "detail",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "post_code",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "longitude",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "latitude",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "is_default",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "address",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpUserBill = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpUserBill",
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
				FieldName: "mp_uid",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "mp_order_uid",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "pm",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "category",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "type",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "number",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "balance",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "mark",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpProductCategory = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpProductCategory",
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
				FieldName: "parent_id",
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
			{
				FieldName: "status",
				Type:      "uint32",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreProduct = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreProduct",
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
			{
				FieldName: "keyword",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "bar_code",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "cate_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "brand_id",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "vip_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "ot_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "postage",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "unit_name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "sort",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "sales",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "stock",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "description",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "give_integral",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "cost",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "ficti",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "browse",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "code_path",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "temp_id",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "spec_type",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "integral",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "prop",
				Type:      "uint32",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreProductAttr = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreProductAttr",
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
				FieldName: "product_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "attr_name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "attr_values",
				Type:      "string",
				IsArray:   true,
			},
		},
	},
}

var ModelMpStoreProductAttrResult = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreProductAttrResult",
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
				FieldName: "product_id",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "result",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreProductAttrValue = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreProductAttrValue",
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
				FieldName: "product_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "sku",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "stock",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "sales",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "image",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "cost",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "bar_code",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "ot_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "weight",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "volume",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "brokerage",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "brokerage_two",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "pink_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "pink_stock",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "seckill_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "seckill_stock",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "integral",
				Type:      "int32",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreProductReply = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreProductReply",
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
				FieldName: "mp_uid",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "mp_order_id",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "product_id",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "reply_type",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "product_score",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "service_score",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "comment",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "pics",
				Type:      "string",
				IsArray:   true,
			},
			{
				FieldName: "merchant_reply_content",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "merchant_reply_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "is_reply",
				Type:      "int32",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreProductRule = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreProductRule",
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
				FieldName: "rule_name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "rule_value",
				Type:      "string",
				IsArray:   true,
			},
		},
	},
}

var ModelMpStoreShop = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreShop",
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
				FieldName: "address_map",
				Type:      "string",
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
			{
				FieldName: "distance",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "min_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "delivery_price",
				Type:      "int64",
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
				FieldName: "admin_id",
				Type:      "string",
				IsArray:   true,
			},
			{
				FieldName: "uniprint_id",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "start_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "end_at",
				Type:      "uint32",
				IsArray:   false,
			},
		},
	},
}

var ModelMpCoupon = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpCoupon",
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
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "title",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "is_switch",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "least",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "value",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "start_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "end_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "weigh",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "type",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "exchange_code",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "receive",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "distribute",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "score",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "instructions",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "image",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "limit",
				Type:      "int32",
				IsArray:   false,
			},
		},
	},
}

var ModelMpCouponUser = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpCouponUser",
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
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "mp_uid",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "coupon_id",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "status",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "exchange_code",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpOrderNumber = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpOrderNumber",
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

var ModelMpStoreOrder = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreOrder",
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
				FieldName: "extend_order_sn",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "mp_uid",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "real_name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "user_phone",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "user_address",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "cart_id",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "freight_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "total_num",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "total_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "total_postage",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "pay_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "pay_postage",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "deduction_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "coupon_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "coupon_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "paid",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "pay_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "pay_type",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "order_type",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "status",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "refund_status",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "refund_reason_wap_img",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "refund_reason_wap_explain",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "refund_reason_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "refund_reason_wap",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "refund_reason",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "refund_price",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "delivery_sn",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "delivery_name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "delivery_type",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "delivery_id",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "delivery_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "gain_integral",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "use_integral",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "pay_integral",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "back_integral",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "mark",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "unique",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "remark",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "mer_id",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "combination_id",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "pink_id",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "cost",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "seckill_id",
				Type:      "int64",
				IsArray:   false,
			},
			{
				FieldName: "bargain_id",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "verify_code",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "mp_store_shop_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "shipping_type",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "is_channel",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "is_system_del",
				Type:      "int32",
				IsArray:   false,
			},
			{
				FieldName: "get_at",
				Type:      "uint32",
				IsArray:   false,
			},
			{
				FieldName: "number_id",
				Type:      "int64",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreOrderCartInfo = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreOrderCartInfo",
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
				FieldName: "cart_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "product_id",
				Type:      "uint64",
				IsArray:   false,
			},
			{
				FieldName: "cart_info",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "unique",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "is_after_sales",
				Type:      "int32",
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
			{
				FieldName: "spec",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "price",
				Type:      "int64",
				IsArray:   false,
			},
		},
	},
}

var ModelMpStoreOrderStatus = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpStoreOrderStatus",
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

var ModelMpMaterial = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpMaterial",
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
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "group_id",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "name",
				Type:      "string",
				IsArray:   false,
			},
			{
				FieldName: "url",
				Type:      "string",
				IsArray:   false,
			},
		},
	},
}

var ModelMpMaterialGroup = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpMaterialGroup",
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
		},
	},
}

var ModelMpService = &engine.ModelObjectType{
	Name: "lbsingle.ModelMpService",
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
	Name: "lbsingle.ModelMpShopAds",
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
