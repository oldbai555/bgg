// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package file

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"postapocgame/admin-server/internal/logic/file"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileDownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 验证必需参数
		if req.Id == 0 {
			httpx.ErrorCtx(r.Context(), w, errs.New(errs.CodeBadRequest, "文件ID不能为空"))
			return
		}

		l := file.NewFileDownloadLogic(r.Context(), svcCtx)
		fileInfo, filePath, err := l.FileDownload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 设置响应头
		// 设置文件名（使用原始文件名）
		fileName := fileInfo.OriginalName
		// 对文件名进行 URL 编码，确保中文文件名正确显示
		w.Header().Set("Content-Disposition", `attachment; filename="`+fileName+`"; filename*=UTF-8''`+url.PathEscape(fileName))

		// 设置 Content-Type
		if fileInfo.MimeType != "" {
			w.Header().Set("Content-Type", fileInfo.MimeType)
		} else {
			// 根据扩展名推断 Content-Type
			ext := filepath.Ext(fileName)
			if ext != "" {
				w.Header().Set("Content-Type", http.DetectContentType([]byte(ext)))
			}
		}

		// 设置文件大小（从文件系统获取实际大小）
		if fileInfo.Size > 0 {
			w.Header().Set("Content-Length", strconv.FormatUint(fileInfo.Size, 10))
		}

		// 直接返回文件内容
		http.ServeFile(w, r, filePath)
	}
}
