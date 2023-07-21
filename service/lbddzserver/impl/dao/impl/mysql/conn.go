package mysql

import (
	"github.com/oldbai555/bgg/pkg/webtool"
	mysql "github.com/oldbai555/driver-mysql"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/gorm/logger"
	"github.com/oldbai555/lbtool/log"
	"time"
)

const (
	autoMigrateOptKey   = "gorm:table_options"
	autoMigrateOptValue = "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"
)

var MasterOrm *gorm.DB
var modelList []interface{}

func RegisterModel(vs ...interface{}) {
	modelList = append(modelList, vs...)
}

func InitMasterOrm(dsn string) error {
	var err error
	MasterOrm, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: gorm.NamingStrategy{
			SingularTable: true,  // 是否单表，命名是否复数
			NoLowerCase:   false, // 是否关闭驼峰命名
		},

		NowFunc: func() int32 {
			return int32(time.Now().Unix())
		},

		PrepareStmt: true, // 预编译 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率

		Logger: webtool.NewOrmLog( //  日志配制
			log.GetLogger(),
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // 日志级别
			}),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func AutoMigrate() {
	if len(modelList) > 0 {
		err := MasterOrm.Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(modelList...)
		if err != nil {
			log.Errorf("err:%v", err)
			panic(err)
		}
	}
}
