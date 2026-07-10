package misc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DailyShortSentenceModel = (*customDailyShortSentenceModel)(nil)

type (
	// DailyShortSentenceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDailyShortSentenceModel.
	DailyShortSentenceModel interface {
		dailyShortSentenceModel
		// WithSession 返回一个绑定到事务 session 的新 DailyShortSentenceModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) DailyShortSentenceModel
	}

	customDailyShortSentenceModel struct {
		*defaultDailyShortSentenceModel
	}
)

// NewDailyShortSentenceModel returns a model for the database table.
func NewDailyShortSentenceModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) DailyShortSentenceModel {
	return &customDailyShortSentenceModel{
		defaultDailyShortSentenceModel: newDailyShortSentenceModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customDailyShortSentenceModel) WithSession(session sqlx.Session) DailyShortSentenceModel {
	return &customDailyShortSentenceModel{
		defaultDailyShortSentenceModel: &defaultDailyShortSentenceModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
