package fitness

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FitnessBodyRecordModel = (*customFitnessBodyRecordModel)(nil)

type (
	// FitnessBodyRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFitnessBodyRecordModel.
	FitnessBodyRecordModel interface {
		fitnessBodyRecordModel
		// WithSession 返回一个绑定到事务 session 的新 FitnessBodyRecordModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) FitnessBodyRecordModel
	}

	customFitnessBodyRecordModel struct {
		*defaultFitnessBodyRecordModel
	}
)

// NewFitnessBodyRecordModel returns a model for the database table.
func NewFitnessBodyRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FitnessBodyRecordModel {
	return &customFitnessBodyRecordModel{
		defaultFitnessBodyRecordModel: newFitnessBodyRecordModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customFitnessBodyRecordModel) WithSession(session sqlx.Session) FitnessBodyRecordModel {
	return &customFitnessBodyRecordModel{
		defaultFitnessBodyRecordModel: &defaultFitnessBodyRecordModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
