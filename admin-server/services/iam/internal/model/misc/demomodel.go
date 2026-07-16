package misc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DemoModel = (*customDemoModel)(nil)

type (
	// DemoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDemoModel.
	DemoModel interface {
		demoModel
		// WithSession 返回一个绑定到事务 session 的新 DemoModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) DemoModel
	}

	customDemoModel struct {
		*defaultDemoModel
	}
)

// NewDemoModel returns a model for the database table.
func NewDemoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) DemoModel {
	return &customDemoModel{
		defaultDemoModel: newDemoModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customDemoModel) WithSession(session sqlx.Session) DemoModel {
	return &customDemoModel{
		defaultDemoModel: &defaultDemoModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
