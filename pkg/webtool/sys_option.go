package webtool

import (
	"fmt"
)

type Option func(tool *WebTool)

func OptionWithOrm(dto ...interface{}) Option {
	return func(tool *WebTool) {
		gorm := &GormMysqlConf{}
		err := gorm.InitConf(tool.V)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
		err = gorm.GenConfTool(tool, dto...)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
	}
}

func OptionWithRdb() Option {
	return func(tool *WebTool) {
		rdb := &RedisConf{}
		err := rdb.InitConf(tool.V)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
		err = rdb.GenConfTool(tool)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
	}
}

func OptionWithStorage() Option {
	return func(tool *WebTool) {
		rdb := &StorageConf{}
		err := rdb.InitConf(tool.V)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
		err = rdb.GenConfTool(tool)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
	}
}

func OptionWithServer() Option {
	return func(tool *WebTool) {
		sc := &ServerConf{}
		err := sc.InitConf(tool.V)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
		err = sc.GenConfTool(tool)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
	}
}
