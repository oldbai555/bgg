package fitness

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FitnessCheckinModel = (*customFitnessCheckinModel)(nil)

type (
	// FitnessCheckinModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFitnessCheckinModel.
	FitnessCheckinModel interface {
		fitnessCheckinModel
		// WithSession 返回一个绑定到事务 session 的新 FitnessCheckinModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) FitnessCheckinModel
	}

	customFitnessCheckinModel struct {
		*defaultFitnessCheckinModel
	}
)

// NewFitnessCheckinModel returns a model for the database table.
func NewFitnessCheckinModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FitnessCheckinModel {
	return &customFitnessCheckinModel{
		defaultFitnessCheckinModel: newFitnessCheckinModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customFitnessCheckinModel) WithSession(session sqlx.Session) FitnessCheckinModel {
	return &customFitnessCheckinModel{
		defaultFitnessCheckinModel: &defaultFitnessCheckinModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
