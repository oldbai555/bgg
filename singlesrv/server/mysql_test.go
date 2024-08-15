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
