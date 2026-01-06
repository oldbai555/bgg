// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sdk

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkCallLogExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkCallLogExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkCallLogExportLogic {
	return &SdkCallLogExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkCallLogExportLogic) SdkCallLogExport(w http.ResponseWriter, r *http.Request, req *types.SdkCallLogExportReq) error {
	if req == nil {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	repo := repository.NewSdkAdminRepository(l.svcCtx.Repository)
	list, err := repo.ExportCallLogs(l.ctx, 2000, req.SdkKeyId, req.ApiCode, req.RespCode, req.Ip, req.StartTime, req.EndTime)
	if err != nil {
		return errs.Wrap(errs.CodeInternalError, "导出调用记录失败", err)
	}

	filename := fmt.Sprintf("sdk_call_log_%s.csv", time.Now().Format("20060102_150405"))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	// BOM for Excel
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(w)
	defer writer.Flush()

	headers := []string{"ID", "SDK Key ID", "接口ID", "API Code", "路径", "方法", "IP", "状态码", "耗时(ms)", "创建时间"}
	if err := writer.Write(headers); err != nil {
		return errs.Wrap(errs.CodeInternalError, "写入CSV表头失败", err)
	}

	for _, log := range list {
		row := []string{
			fmt.Sprintf("%d", log.Id),
			fmt.Sprintf("%d", log.SdkKeyId),
			fmt.Sprintf("%d", log.SdkInterfaceId),
			log.ApiCode,
			log.Path,
			log.Method,
			log.Ip,
			fmt.Sprintf("%d", log.RespCode),
			fmt.Sprintf("%d", log.DurationMs),
			time.Unix(log.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return errs.Wrap(errs.CodeInternalError, "写入CSV数据失败", err)
		}
	}

	return nil
}
