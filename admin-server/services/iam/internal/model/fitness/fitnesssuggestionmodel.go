package fitness

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FitnessSuggestionModel = (*customFitnessSuggestionModel)(nil)

type (
	// FitnessSuggestionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFitnessSuggestionModel.
	FitnessSuggestionModel interface {
		fitnessSuggestionModel
		// WithSession 返回一个绑定到事务 session 的新 FitnessSuggestionModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) FitnessSuggestionModel
	}

	customFitnessSuggestionModel struct {
		*defaultFitnessSuggestionModel
	}
)

// NewFitnessSuggestionModel returns a model for the database table.
func NewFitnessSuggestionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FitnessSuggestionModel {
	return &customFitnessSuggestionModel{
		defaultFitnessSuggestionModel: newFitnessSuggestionModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customFitnessSuggestionModel) WithSession(session sqlx.Session) FitnessSuggestionModel {
	return &customFitnessSuggestionModel{
		defaultFitnessSuggestionModel: &defaultFitnessSuggestionModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
