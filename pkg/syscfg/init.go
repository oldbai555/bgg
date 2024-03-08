package syscfg

var Global *SysCfg

func InitGlobal(srv, path string, opt ...Option) {
	Global = LoadSysCfgByYaml(srv, path, opt...)
}
