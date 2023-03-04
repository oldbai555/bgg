package impl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/client/lbblog"
	"github.com/oldbai555/bgg/client/lbuser"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/spf13/viper"

	"github.com/oldbai555/lbtool/log"
	"github.com/storyicon/grbac"
	"time"
)

var lb *Tool

type Tool struct {
	*webtool.WebTool
	Rbac       *grbac.Controller
	WechatConf WeChatGzhConf
}

func StartServer() error {
	v, err := initViper()
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	lb = &Tool{}
	lb.WebTool, err = webtool.NewWebTool(v, webtool.OptionWithOrm(
		&lbuser.ModelUser{},
		&lbblog.ModelArticle{},
		&lbblog.ModelCategory{},
		&lbblog.ModelComment{},
	), webtool.OptionWithRdb(), webtool.OptionWithStorage())

	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	val := lb.V.Get("wechatConf")
	err = webtool.JsonConvertStruct(val, &lb.WechatConf)
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}

	gin.DefaultWriter = log.GetWriter()
	h := gin.New()

	// 配置跨域中间件
	// 链路追踪
	h.Use(Cors(), RegisterUuidTrace(), gin.LoggerWithFormatter(defaultLogFormatter), gin.Recovery(), RegisterShowReq())

	// 权限管理
	// 在这里，我们通过“grbac.WithLoader”接口使用自定义Loader功能
	// 并指定应每分钟调用一次LoadAuthorizationRules函数以获取最新的身份验证规则。
	if lb.Rbac, err = grbac.New(grbac.WithLoader(LoadAuthorizationRules, time.Minute)); err != nil {
		log.Errorf("err is : %v", err)
		return err
	}

	// 初始化接口
	InitSysApi(h)

	// 初始化数据库工具类
	InitDbOrm()

	// 启动服务
	err = h.Run(fmt.Sprintf(":%d", lb.Sc.Port))
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func InitSysApi(h *gin.Engine) {
	registerLbuserApi(h)
	registerLbblogApi(h)
	registerStoreApi(h)
	registerPublicApi(h)
	registerLbwebsocketApi(h)
	registerWechatGzhApi(h)
}

func initViper() (*viper.Viper, error) {
	viper.SetConfigName("application")          // name of config file (without extension)
	viper.SetConfigType("yaml")                 // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/work/")           // path to look for the config file in
	viper.AddConfigPath("./")                   // optionally look for config in the working directory
	viper.AddConfigPath("./lbserver/cmd/")      // optionally look for config in the working directory
	viper.AddConfigPath("../cmd/")              // optionally look for config in the working directory
	viper.AddConfigPath("./lbserver/resource/") // optionally look for config in the working directory
	viper.AddConfigPath("../resource/")         // optionally look for config in the working directory
	err := viper.ReadInConfig()                 // Find and read the config file
	if err != nil {                             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper(), nil
}
