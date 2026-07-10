package video

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VideoModel = (*customVideoModel)(nil)

type (
	// VideoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideoModel.
	VideoModel interface {
		videoModel
		// WithSession 返回一个绑定到事务 session 的新 VideoModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) VideoModel
	}

	customVideoModel struct {
		*defaultVideoModel
	}
)

// NewVideoModel returns a model for the database table.
func NewVideoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) VideoModel {
	return &customVideoModel{
		defaultVideoModel: newVideoModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customVideoModel) WithSession(session sqlx.Session) VideoModel {
	return &customVideoModel{
		defaultVideoModel: &defaultVideoModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
