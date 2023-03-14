package webtool

type Option func(tool *WebTool)

func OptionWithOrm() Option {
	return func(tool *WebTool) {
		tool.GormMysqlConf = NewGormMysqlConf(tool.V)
	}
}

func OptionWithRdb() Option {
	return func(tool *WebTool) {
		tool.RedisConf = NewRedisConf(tool.V)
	}
}

func OptionWithStorage() Option {
	return func(tool *WebTool) {
		tool.StorageConf = NewStorageConf(tool.V)
	}
}

func OptionWithServer() Option {
	return func(tool *WebTool) {
		tool.ServerConf = NewServerConf(tool.V)
	}
}

func OptionWithWxGzh() Option {
	return func(tool *WebTool) {
		tool.WxGzhConf = NewWxGzhConf(tool.V)
	}
}
