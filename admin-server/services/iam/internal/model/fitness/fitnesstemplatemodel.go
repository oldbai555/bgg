package fitness

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FitnessTemplateModel = (*customFitnessTemplateModel)(nil)

type (
	// FitnessTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFitnessTemplateModel.
	FitnessTemplateModel interface {
		fitnessTemplateModel
		// WithSession 返回一个绑定到事务 session 的新 FitnessTemplateModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) FitnessTemplateModel
	}

	customFitnessTemplateModel struct {
		*defaultFitnessTemplateModel
	}
)

// NewFitnessTemplateModel returns a model for the database table.
func NewFitnessTemplateModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FitnessTemplateModel {
	return &customFitnessTemplateModel{
		defaultFitnessTemplateModel: newFitnessTemplateModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customFitnessTemplateModel) WithSession(session sqlx.Session) FitnessTemplateModel {
	return &customFitnessTemplateModel{
		defaultFitnessTemplateModel: &defaultFitnessTemplateModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
