package system

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminDictItemModel = (*customAdminDictItemModel)(nil)

type (
	// AdminDictItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminDictItemModel.
	AdminDictItemModel interface {
		adminDictItemModel
		FindPageByTypeId(ctx context.Context, typeId uint64, page, pageSize int64) ([]AdminDictItem, int64, error)
	}

	customAdminDictItemModel struct {
		*defaultAdminDictItemModel
	}
)

// NewAdminDictItemModel returns a model for the database table.
func NewAdminDictItemModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminDictItemModel {
	return &customAdminDictItemModel{
		defaultAdminDictItemModel: newAdminDictItemModel(conn, c, opts...),
	}
}

// FindPageByTypeId 按 typeId 分页查询字典项
func (m *customAdminDictItemModel) FindPageByTypeId(ctx context.Context, typeId uint64, page, pageSize int64) ([]AdminDictItem, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize

	where := "where `type_id` = ? and deleted_at = 0"

	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, where)
	if err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, typeId); err != nil {
		return nil, 0, err
	}

	var list []AdminDictItem
	query := fmt.Sprintf("select %s from %s %s order by sort asc, id asc limit ? offset ?", adminDictItemRows, m.table, where)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, typeId, pageSize, offset); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
