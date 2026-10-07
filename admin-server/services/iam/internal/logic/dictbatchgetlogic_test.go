package logic

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"postapocgame/admin-server/services/iam/iam"
)

var (
	dictTypeColumns = []string{"id", "name", "code", "description", "status", "created_at", "updated_at", "deleted_at"}
	dictItemColumns = []string{"id", "type_id", "label", "value", "sort", "status", "remark", "created_at", "updated_at", "deleted_at"}
)

func TestDictBatchGet_QueriesAtMostTwice(t *testing.T) {
	svcCtx, sqlMock, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	sqlMock.ExpectQuery("SELECT \\* FROM admin_dict_type WHERE").
		WillReturnRows(sqlmock.NewRows(dictTypeColumns).
			AddRow(1, "类型A", "type_a", "", 1, 0, 0, 0).
			AddRow(2, "类型B", "type_b", "", 1, 0, 0, 0))
	sqlMock.ExpectQuery("SELECT \\* FROM admin_dict_item WHERE").
		WillReturnRows(sqlmock.NewRows(dictItemColumns).
			AddRow(11, 1, "甲", "1", 1, 1, nil, 0, 0, 0).
			AddRow(12, 1, "乙", "2", 2, 1, "备注", 0, 0, 0))

	resp, err := NewDictBatchGetLogic(context.Background(), svcCtx).DictBatchGet(&iam.DictBatchGetRequest{
		Codes: []string{"type_a", "type_b", "missing", "type_a", ""},
	})
	require.NoError(t, err)
	require.NoError(t, sqlMock.ExpectationsWereMet())

	require.Len(t, resp.Dicts, 2)
	a := resp.Dicts["type_a"]
	require.Len(t, a.Items, 2)
	assert.Equal(t, "甲", a.Items[0].Label)
	assert.Equal(t, "备注", a.Items[1].Remark)
	assert.Empty(t, resp.Dicts["type_b"].Items)
	assert.NotContains(t, resp.Dicts, "missing")
}

func TestDictBatchGet_CachedTypeSkipsItemQuery(t *testing.T) {
	svcCtx, sqlMock, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	require.NoError(t, svcCtx.Repository.BusinessCache.SetDictItems(context.Background(), "type_a",
		[]*iam.DictItemItem{{Id: 11, TypeId: 1, Label: "缓存", Value: "1"}}))

	sqlMock.ExpectQuery("SELECT \\* FROM admin_dict_type WHERE").
		WillReturnRows(sqlmock.NewRows(dictTypeColumns).AddRow(1, "类型A", "type_a", "", 1, 0, 0, 0))

	resp, err := NewDictBatchGetLogic(context.Background(), svcCtx).DictBatchGet(&iam.DictBatchGetRequest{
		Codes: []string{"type_a"},
	})
	require.NoError(t, err)
	require.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Equal(t, "缓存", resp.Dicts["type_a"].Items[0].Label)
}

func TestDictBatchGet_EmptyCodes(t *testing.T) {
	svcCtx, _, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	_, err := NewDictBatchGetLogic(context.Background(), svcCtx).DictBatchGet(&iam.DictBatchGetRequest{})
	require.Error(t, err)

	resp, err := NewDictBatchGetLogic(context.Background(), svcCtx).DictBatchGet(&iam.DictBatchGetRequest{Codes: []string{""}})
	require.NoError(t, err)
	assert.Empty(t, resp.Dicts)
}
