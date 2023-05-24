package mysql

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/pkg/webtool"
	mysql "github.com/oldbai555/driver-mysql"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/gorm/logger"
	"github.com/oldbai555/lbtool/log"
	"sync"
	"time"
)

const (
	autoMigrateOptKey   = "gorm:table_options"
	autoMigrateOptValue = "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"
)

type mysqlConn struct {
	conn   *gorm.DB
	dsn    string
	connMu sync.Mutex
}

func (c *mysqlConn) mustGetConn(ctx context.Context) *gorm.DB {
	if conn, err := c.getConn(ctx); err != nil {
		panic(any(err))
	} else {
		return conn
	}
}

func (c *mysqlConn) getConn(ctx context.Context) (*gorm.DB, error) {
	if c.conn != nil {
		return c.conn.WithContext(ctx), nil
	}
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.conn != nil {
		return c.conn.WithContext(ctx), nil
	}
	var err error
	c.conn, err = gorm.Open(mysql.Open(c.dsn), &gorm.Config{
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
		log.Errorf(fmt.Sprintf("err is : %v", err))
		return nil, err
	}

	return c.conn.WithContext(ctx), nil
}
