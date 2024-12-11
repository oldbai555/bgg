/**
 * @Author: zjj
 * @Date: 2024/6/18
 * @Desc:
**/

package lbsingleserver

import (
	"fmt"
	"github.com/oldbai555/bgg/pkg/syscfg"
	client2 "github.com/oldbai555/bgg/service/lbsingle/client"
	"github.com/oldbai555/bgg/service/lbsingle/client/lbsingledb"
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
	gormEngine.AutoMigrate([]interface{}{&client2.ModelFile{}})
	gormEngine.RegObjectType(lbsingledb.ModelFile)
	OrmFile := gormx.NewBaseModel[*client2.ModelFile](gormx.ModelConfig{
		NotFoundErrCode: int32(client2.ErrCode_ErrFileNotFound),
		Db:              "biz",
	})
	find, err := OrmFile.NewBaseScope().Where(client2.FieldId_, 3).First(uctx.NewBaseUCtx())
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	var file = &client2.ModelFile{
		Size:    1,
		Name:    "2",
		Rename:  "3",
		Path:    "4",
		Md5:     utils.GenUUID(),
		SortUrl: utils.GenUUID(),
		State:   0,
	}
	var file1 = &client2.ModelFile{
		Size:    1,
		Name:    "2",
		Rename:  "3",
		Path:    "4",
		Md5:     utils.GenUUID(),
		SortUrl: utils.GenUUID(),
		State:   0,
	}
	_, err = OrmFile.NewBaseScope().BatchCreate(uctx.NewBaseUCtx(), 0, []*client2.ModelFile{
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
	var user = &client2.ModelUser{
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

func TestAddService(t *testing.T) {

}
