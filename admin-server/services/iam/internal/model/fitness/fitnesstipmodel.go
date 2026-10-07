package fitness

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FitnessTipModel = (*customFitnessTipModel)(nil)

type (
	// FitnessTipModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFitnessTipModel.
	FitnessTipModel interface {
		fitnessTipModel
		// WithSession 返回一个绑定到事务 session 的新 FitnessTipModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) FitnessTipModel
	}

	customFitnessTipModel struct {
		*defaultFitnessTipModel
	}
)

// NewFitnessTipModel returns a model for the database table.
func NewFitnessTipModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FitnessTipModel {
	return &customFitnessTipModel{
		defaultFitnessTipModel: newFitnessTipModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customFitnessTipModel) WithSession(session sqlx.Session) FitnessTipModel {
	return &customFitnessTipModel{
		defaultFitnessTipModel: &defaultFitnessTipModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
