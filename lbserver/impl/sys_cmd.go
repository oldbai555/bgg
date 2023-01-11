package impl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/lbblog"
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/bgg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/storyicon/grbac"
	"time"
)

const defaultApolloServerPrefix = "server"

var lb *Tool

type Tool struct {
	*webtool.WebTool
	Rbac *grbac.Controller
	Sc   ServerConfig
}

type ServerConfig struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

func StartServer() {
	var err error
	lb = &Tool{}
	lb.WebTool, err = webtool.NewWebTool(&webtool.ApolloConf{
		AppId:     "golb",
		NameSpace: "application.yaml",
		Address:   "http://127.0.0.1:8080",
		Cluster:   "DEV",
		Secret:    "0e0c07411823424a9aeef49b20d046ce",
	}, webtool.OptionWithOrm(
		&lbuser.ModelUser{},
		&lbblog.ModelArticle{},
		&lbblog.ModelCategory{},
		&lbblog.ModelComment{},
	), webtool.OptionWithRdb(),
		webtool.OptionWithStorage())

	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// 拿到服务的信息
	err = lb.GetJson4Apollo(defaultApolloServerPrefix, &lb.Sc)
	if err != nil {
		log.Errorf("err is : %v", err)
		return
	}
	log.Infof("init server successfully")

	gin.DefaultWriter = log.GetWriter()
	h := gin.New()

	// 配置跨域中间件
	// 链路追踪
	h.Use(Cors(), RegisterUuidTrace(), gin.LoggerWithFormatter(defaultLogFormatter), gin.Recovery())

	// 权限管理
	// 在这里，我们通过“grbac.WithLoader”接口使用自定义Loader功能
	// 并指定应每分钟调用一次LoadAuthorizationRules函数以获取最新的身份验证规则。
	if lb.Rbac, err = grbac.New(grbac.WithLoader(LoadAuthorizationRules, time.Minute)); err != nil {
		log.Errorf("err is : %v", err)
		return
	}

	// 初始化接口
	InitSysApi(h)

	// 初始化数据库工具类
	InitDbOrm()

	// 启动服务
	err = h.Run(fmt.Sprintf(":%d", lb.Sc.Port))
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
}

func InitSysApi(h *gin.Engine) {
	registerLbuserApi(h)
	registerLbblogApi(h)
	registerStoreApi(h)
}
