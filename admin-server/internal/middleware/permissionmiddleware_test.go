package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	iamdomain "postapocgame/admin-server/internal/domain/iam"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
)

// 08-testing-strategy.md §4：PermissionResolver.CanAccess 的多分支逻辑已经在
// internal/domain/iam/permission_resolver_test.go 覆盖，这里只做“权限不足时返回正确错误码”
// 的中间件层烟测，不重复分支组合。
func TestPermissionMiddleware_Handle_Forbidden(t *testing.T) {
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	defer db.Close()
	conn := sqlx.NewSqlConnFromDB(db)

	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	redisConf := redis.RedisConf{Host: mr.Addr(), Type: "node"}
	cacheConf := cache.CacheConf{{RedisConf: redisConf, Weight: 100}}
	rdb, err := redis.NewRedis(redisConf)
	require.NoError(t, err)

	repo, err := repository.NewRepository(conn, cacheConf, rdb)
	require.NoError(t, err)

	// 非超级管理员（userID=2）访问一个存在但用户没有任何角色的接口，CanAccess 分支
	// 走到 FindByMethodAndPath 命中、ListRoleIDsByUserID 返回空 -> (false, nil)。
	sqlMock.ExpectQuery(`from ` + "`" + `admin_api` + "`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "method", "path", "description", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(10, "用户列表", "GET", "/api/v1/iam/user/list", "", 1, 0, 0, 0))
	sqlMock.ExpectQuery("from admin_user_role").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "role_id", "created_at", "updated_at"}))

	resolver := iamdomain.NewPermissionResolver(repo)
	m := NewPermissionMiddleware(resolver)

	called := false
	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/iam/user/list", nil)
	req = req.WithContext(jwthelper.WithAuthUser(req.Context(), jwthelper.AuthUser{UserID: 2, Username: "bob"}))
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.False(t, called, "无权限时不应该放行到下一个 handler")
	assert.Equal(t, http.StatusBadRequest, rec.Code) // pkg/response.ErrorCtx 对业务错误统一写 400

	var body struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, errs.CodeForbidden, body.Code)
}

func TestPermissionMiddleware_Handle_Unauthenticated(t *testing.T) {
	m := NewPermissionMiddleware(iamdomain.NewPermissionResolver(nil))

	called := false
	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/iam/user/list", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var body struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, errs.CodeUnauthorized, body.Code)
}
