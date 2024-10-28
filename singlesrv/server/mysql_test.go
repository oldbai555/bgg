/**
 * @Author: zjj
 * @Date: 2024/6/18
 * @Desc:
**/

package server

import (
	"fmt"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/client/lbsingledb"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/gormx/egimpl"
	"github.com/oldbai555/micro/gormx/engine"
	"github.com/oldbai555/micro/uctx"
	"testing"
)

func initMysqlTest() {
	mysqlConf := syscfg.NewGormMysqlConf("")
	gormEngine := egimpl.NewGormEngine(mysqlConf.Dsn())
	engine.SetOrmEngine(gormEngine)
	Init()
}

func TestInit(t *testing.T) {
	mysqlConf := syscfg.NewGormMysqlConf("")
	gormEngine := egimpl.NewGormEngine(mysqlConf.Dsn())
	engine.SetOrmEngine(gormEngine)
	gormEngine.AutoMigrate([]interface{}{&client.ModelFile{}})
	gormEngine.RegObjectType(lbsingledb.ModelFile)
	OrmFile := gormx.NewBaseModel[*client.ModelFile](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrFileNotFound),
		Db:              "biz",
	})
	find, err := OrmFile.NewBaseScope().Where(client.FieldId_, 3).First(uctx.NewBaseUCtx())
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	var file = &client.ModelFile{
		Size:    1,
		Name:    "2",
		Rename:  "3",
		Path:    "4",
		Md5:     utils.GenUUID(),
		SortUrl: utils.GenUUID(),
		State:   0,
	}
	var file1 = &client.ModelFile{
		Size:    1,
		Name:    "2",
		Rename:  "3",
		Path:    "4",
		Md5:     utils.GenUUID(),
		SortUrl: utils.GenUUID(),
		State:   0,
	}
	_, err = OrmFile.NewBaseScope().BatchCreate(uctx.NewBaseUCtx(), 0, []*client.ModelFile{
		file, file1,
	})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	OrmFile.NewBaseScope().Where("id", file.Id).Update(uctx.NewBaseUCtx(), map[string]interface{}{
		"name": "修改数据" + fmt.Sprintf("%d", utils.TimeNow()),
	})
	file1.Name = "1修改数据" + fmt.Sprintf("%d", utils.TimeNow())
	OrmFile.NewBaseScope().Save(uctx.NewBaseUCtx(), file1)
	OrmFile.NewBaseScope().Where("id", file.Id).Delete(uctx.NewBaseUCtx())

	t.Log(find)
	t.Log(file)
	t.Log(file1)
}
func TestAddUser(t *testing.T) {
	initMysqlTest()
	var user = &client.ModelUser{
		Username: "bigbai003",
		Password: "123456",
		Nickname: "大白3号",
	}

	err := OrmUser.NewBaseScope().Create(uctx.NewBaseUCtx(), user)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
func TestDecStock(t *testing.T) {
	initMysqlTest()
	_, err := OrmMpStoreProductAttrValue.NewBaseScope().Where(map[string]interface{}{
		client.FieldProductId_: 1,
		client.FieldSku_:       "小条",
	}).Update(uctx.NewBaseUCtx(), map[string]interface{}{
		client.FieldStock_: fmt.Sprintf("%s-1", client.FieldStock_),
	})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}

func TestAddStore(t *testing.T) {
	initMysqlTest()
	var store = &client.ModelMpStoreShop{
		Name:          "大白的食谱",
		Mobile:        "18877227897",
		Image:         "https://oldbai.top/oss/download/BUOZ74",
		Images:        []string{"https://oldbai.top/oss/download/BUOZ74", "https://oldbai.top/oss/download/BUOZ74"},
		Address:       "广东省广州市海珠区官洲街道",
		AddressMap:    "",
		Distance:      0,
		MinPrice:      0,
		DeliveryPrice: 0,
		Notice:        "选择喜欢的食物吧~",
		Status:        1,
		AdminId:       []string{"1"},
		UniprintId:    "",
		StartAt:       1729842801,
		EndAt:         1730274801,
	}
	err := OrmMpStoreShop.NewBaseScope().Create(uctx.NewBaseUCtx(), store)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	var cate = &client.ModelMpProductCategory{
		MpStoreShopId: store.Id,
		Name:          "水产",
		PicUrl:        "https://oldbai.top/oss/download/BUOZ74",
		Description:   "鲜活水产",
	}
	err = OrmMpProductCategory.NewBaseScope().Create(uctx.NewBaseUCtx(), cate)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	var prop = uint32(0)
	prop = utils.SetBit(prop, uint32(client.ModelMpStoreProduct_PropShow))
	prop = utils.SetBit(prop, uint32(client.ModelMpStoreProduct_PropHot))
	prop = utils.SetBit(prop, uint32(client.ModelMpStoreProduct_PropIntegral))
	var product = &client.ModelMpStoreProduct{
		MpStoreShopId: store.Id,
		Image:         "https://oldbai.top/oss/download/BUOZ74",
		SliderImage:   []string{"https://oldbai.top/oss/download/BUOZ74", "https://oldbai.top/oss/download/BUOZ74"},
		Name:          "清蒸鲈鱼",
		Info:          "清蒸鲈鱼,鲜美",
		Keyword:       "鲈鱼",
		CateId:        cate.Id,
		Price:         188,
		VipPrice:      168,
		OtPrice:       188,
		UnitName:      "条",
		Sort:          1,
		Sales:         5,
		Stock:         1,
		Description:   "一条鲜美的鲈鱼配上姜丝、葱丝，热油淋上可美味了",
		GiveIntegral:  188,
		Cost:          150,
		Ficti:         20,
		Browse:        10,
		Integral:      288,
		Prop:          prop,
		SpecType:      1,
	}
	err = OrmMpStoreProduct.NewBaseScope().Create(uctx.NewBaseUCtx(), product)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	var rule = &client.ModelMpStoreProductRule{
		RuleName:  "块头",
		RuleValue: []string{"小条", "中条", "大条"},
	}
	err = OrmMpStoreProductRule.NewBaseScope().Create(uctx.NewBaseUCtx(), rule)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	var attr = &client.ModelMpStoreProductAttr{
		ProductId:  product.Id,
		AttrName:   "块头",
		AttrValues: []string{"小条", "中条", "大条"},
	}

	err = OrmMpStoreProductAttr.NewBaseScope().Create(uctx.NewBaseUCtx(), attr)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	var attrValue1 = &client.ModelMpStoreProductAttrValue{
		ProductId:    product.Id,
		Sku:          "小条",
		Stock:        20,
		Sales:        10,
		Price:        198,
		Image:        "https://oldbai.top/oss/download/BUOZ74",
		Cost:         188,
		OtPrice:      208,
		Weight:       180,
		Volume:       300,
		PinkPrice:    178,
		PinkStock:    10,
		SeckillPrice: 168,
		SeckillStock: 5,
		Integral:     268,
	}
	var attrValue2 = &client.ModelMpStoreProductAttrValue{
		ProductId:    product.Id,
		Sku:          "中条",
		Stock:        30,
		Sales:        20,
		Price:        208,
		Image:        "https://oldbai.top/oss/download/BUOZ74",
		Cost:         198,
		OtPrice:      218,
		Weight:       200,
		Volume:       320,
		PinkPrice:    188,
		PinkStock:    15,
		SeckillPrice: 178,
		SeckillStock: 7,
		Integral:     278,
	}
	var attrValue3 = &client.ModelMpStoreProductAttrValue{
		ProductId:    product.Id,
		Sku:          "大条",
		Stock:        40,
		Sales:        30,
		Price:        218,
		Image:        "https://oldbai.top/oss/download/BUOZ74",
		Cost:         208,
		OtPrice:      228,
		Weight:       220,
		Volume:       340,
		PinkPrice:    198,
		PinkStock:    20,
		SeckillPrice: 188,
		SeckillStock: 10,
		Integral:     288,
	}
	_, err = OrmMpStoreProductAttrValue.NewBaseScope().BatchCreate(uctx.NewBaseUCtx(), 2000, []*client.ModelMpStoreProductAttrValue{attrValue1, attrValue2, attrValue3})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
