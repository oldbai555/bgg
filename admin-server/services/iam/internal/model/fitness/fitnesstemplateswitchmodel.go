package fitness

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FitnessTemplateSwitchModel = (*customFitnessTemplateSwitchModel)(nil)

type (
	// FitnessTemplateSwitchModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFitnessTemplateSwitchModel.
	FitnessTemplateSwitchModel interface {
		fitnessTemplateSwitchModel
		// WithSession 返回一个绑定到事务 session 的新 FitnessTemplateSwitchModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) FitnessTemplateSwitchModel
	}

	customFitnessTemplateSwitchModel struct {
		*defaultFitnessTemplateSwitchModel
	}
)

// NewFitnessTemplateSwitchModel returns a model for the database table.
func NewFitnessTemplateSwitchModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FitnessTemplateSwitchModel {
	return &customFitnessTemplateSwitchModel{
		defaultFitnessTemplateSwitchModel: newFitnessTemplateSwitchModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customFitnessTemplateSwitchModel) WithSession(session sqlx.Session) FitnessTemplateSwitchModel {
	return &customFitnessTemplateSwitchModel{
		defaultFitnessTemplateSwitchModel: &defaultFitnessTemplateSwitchModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
