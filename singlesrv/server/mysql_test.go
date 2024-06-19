/**
 * @Author: zjj
 * @Date: 2024/6/18
 * @Desc:
**/

package server

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/client/lbsingledb"
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
	find, err := OrmFile.NewBaseScope().Where(client.FieldId_, 1).First(uctx.NewBaseUCtx())
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(find)
}
