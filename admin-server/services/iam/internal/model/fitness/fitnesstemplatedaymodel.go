package fitness

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FitnessTemplateDayModel = (*customFitnessTemplateDayModel)(nil)

type (
	// FitnessTemplateDayModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFitnessTemplateDayModel.
	FitnessTemplateDayModel interface {
		fitnessTemplateDayModel
		// WithSession 返回一个绑定到事务 session 的新 FitnessTemplateDayModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) FitnessTemplateDayModel
	}

	customFitnessTemplateDayModel struct {
		*defaultFitnessTemplateDayModel
	}
)

// NewFitnessTemplateDayModel returns a model for the database table.
func NewFitnessTemplateDayModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FitnessTemplateDayModel {
	return &customFitnessTemplateDayModel{
		defaultFitnessTemplateDayModel: newFitnessTemplateDayModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customFitnessTemplateDayModel) WithSession(session sqlx.Session) FitnessTemplateDayModel {
	return &customFitnessTemplateDayModel{
		defaultFitnessTemplateDayModel: &defaultFitnessTemplateDayModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
