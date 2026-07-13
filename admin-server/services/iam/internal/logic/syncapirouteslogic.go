package logic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncApiRoutesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncApiRoutesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiRoutesLogic {
	return &SyncApiRoutesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SyncApiRoutes 原样迁移自 admin.go 的 syncRoutesToAdminAPI：gateway 启动时把已注册的
// REST 路由同步进 admin_api 表（已存在则跳过），拆分后 gateway 没有直连数据库的能力，
// 改成一次性批量 RPC，逻辑本身不变。
func (l *SyncApiRoutesLogic) SyncApiRoutes(in *iam.SyncApiRoutesRequest) (*iam.Empty, error) {
	now := time.Now().Unix()

	for _, r := range in.Routes {
		method := strings.ToUpper(strings.TrimSpace(r.Method))
		path := strings.TrimSpace(r.Path)
		if method == "" || path == "" {
			continue
		}

		_, err := l.svcCtx.Repository.AdminApiModel.FindOneByMethodPath(l.ctx, method, path)
		if err == nil {
			continue
		}
		if !isErrNotFound(err) {
			return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询接口失败", err))
		}

		apiName := fmt.Sprintf("%s_%s", method, sanitizePathForName(path))
		data := &iammodel.AdminApi{
			Name:        apiName,
			Method:      method,
			Path:        path,
			Description: sql.NullString{},
			Status:      1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if _, err = l.svcCtx.Repository.AdminApiModel.Insert(l.ctx, data); err != nil {
			l.Errorf("写入接口 %s %s 失败: %v", method, path,
				errs.Wrap(errs.CodeInternalError, "同步接口路由失败", err))
		}
	}

	return &iam.Empty{}, nil
}

func sanitizePathForName(path string) string {
	path = strings.Trim(path, "/")
	if path == "" {
		return "ROOT"
	}
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, ":", "_")
	return path
}
