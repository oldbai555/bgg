package fitness

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FitnessProfileModel = (*customFitnessProfileModel)(nil)

type (
	// FitnessProfileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFitnessProfileModel.
	FitnessProfileModel interface {
		fitnessProfileModel
		// WithSession 返回一个绑定到事务 session 的新 FitnessProfileModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) FitnessProfileModel
	}

	customFitnessProfileModel struct {
		*defaultFitnessProfileModel
	}
)

// NewFitnessProfileModel returns a model for the database table.
func NewFitnessProfileModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FitnessProfileModel {
	return &customFitnessProfileModel{
		defaultFitnessProfileModel: newFitnessProfileModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customFitnessProfileModel) WithSession(session sqlx.Session) FitnessProfileModel {
	return &customFitnessProfileModel{
		defaultFitnessProfileModel: &defaultFitnessProfileModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
