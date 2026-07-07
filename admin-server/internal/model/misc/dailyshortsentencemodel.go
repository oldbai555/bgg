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
