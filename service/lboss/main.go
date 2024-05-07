package main

import (
	"context"
	"github.com/oldbai555/bgg/pkg/ginhelper"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"path"
)

var dbConn *leveldb.DB

func main() {
	ctx := context.Background()

	// 初始化levelDB k-v 数据库
	var err error
	dbConn, err = leveldb.OpenFile(path.Join(utils.GetCurDir(), "levelDB"), nil)
	if err != nil {
		log.Warnf("err is %v", err)
		return
	}
	defer func(dbConn *leveldb.DB) {
		err := dbConn.Close()
		if err != nil {
			panic(err)
		}
	}(dbConn)

	// 同步最新的文件索引
	_, err = syncFileIndex()
	if err != nil {
		log.Warnf("err is %v", err)
		return
	}

	// 启动服务
	err = server(ctx)
	if err != nil {
		log.Warnf("err is %v", err)
		return
	}
}

func server(ctx context.Context) error {
	syscfg.InitGlobal("", utils.GetCurDir(), syscfg.OptionWithServer())
	srvName := syscfg.Global.ServerConf.Name
	port := syscfg.Global.ServerConf.Port
	log.SetModuleName(srvName)
	return ginhelper.QuickStart(ctx, srvName, port, registerRouter)
}
