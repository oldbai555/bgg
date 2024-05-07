package syscfg

type Option func(tool *SysCfg)

func OptionWithStorage() Option {
	return func(tool *SysCfg) {
		tool.StorageConf = NewStorageConf(tool.V)
	}
}

func OptionWithServer() Option {
	return func(tool *SysCfg) {
		tool.ServerConf = NewServerConf(tool.V)
	}
}
