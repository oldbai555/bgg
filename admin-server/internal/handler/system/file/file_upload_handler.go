// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/system/file"
	"postapocgame/admin-server/internal/svc"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 文件上传需要处理 multipart/form-data
		l := file.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
