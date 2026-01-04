// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package file

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(r *http.Request) (resp *types.FileUploadResp, err error) {
	// 解析 multipart/form-data
	err = r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadRequest, "解析上传文件失败", err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadRequest, "获取上传文件失败", err)
	}
	defer file.Close()

	// 创建上传目录（如果不存在）
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "创建上传目录失败", err)
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	// 文件系统路径（用于实际存储）
	fileSystemPath := filepath.Join(uploadDir, fileName)
	// 访问路径（相对路径，如 /uploads/xxx，用于拼接 URL）
	accessPath := fmt.Sprintf("/uploads/%s", fileName)
	// 获取基础 URL（从配置中读取）
	baseURL := strings.TrimSuffix(l.svcCtx.Config.BaseURL, "/")

	// 保存文件
	dst, err := os.Create(fileSystemPath)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "创建文件失败", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "保存文件失败", err)
	}

	// 获取文件大小
	fileInfo, err := os.Stat(fileSystemPath)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "获取文件信息失败", err)
	}

	// 获取 MIME 类型
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType([]byte(ext))
	}

	// 保存文件记录到数据库
	fileModel := model.AdminFile{
		Name:         fileName,
		OriginalName: header.Filename,
		Path:         accessPath, // 访问路径（相对路径）
		BaseUrl:      baseURL,    // 基础 URL
		Size:         uint64(fileInfo.Size()),
		MimeType:     sql.NullString{String: mimeType, Valid: mimeType != ""},
		Ext:          sql.NullString{String: strings.TrimPrefix(ext, "."), Valid: ext != ""},
		StorageType:  "local",
		Status:       1,
	}

	fileRepo := repository.NewFileRepository(l.svcCtx.Repository)
	if err := fileRepo.Create(l.ctx, &fileModel); err != nil {
		// 如果数据库保存失败，删除已上传的文件
		os.Remove(fileSystemPath)
		return nil, errs.Wrap(errs.CodeInternalError, "保存文件记录失败", err)
	}

	// 返回nginx代理路径（用于前端访问）
	// nginx配置中，/files/uploads/ 路径会代理到后端的 /api/v1/uploads/
	// 文件访问路径格式：/files/uploads/xxx
	proxyPath := fmt.Sprintf("/files/uploads/%s", fileName)

	// 兼容字段：如果配置了BaseURL，也返回完整URL
	fullURL := proxyPath
	if baseURL != "" {
		// 如果BaseURL包含域名，使用完整URL；否则只返回路径
		if strings.HasPrefix(baseURL, "http://") || strings.HasPrefix(baseURL, "https://") {
			fullURL = fmt.Sprintf("%s%s", baseURL, proxyPath)
		} else {
			fullURL = proxyPath
		}
	}

	return &types.FileUploadResp{
		Id:           fileModel.Id,
		Name:         fileModel.Name,
		OriginalName: fileModel.OriginalName,
		Path:         proxyPath, // 使用nginx代理路径
		BaseUrl:      "",        // nginx代理不需要BaseURL
		Url:          fullURL,   // 兼容字段，返回完整 URL
		Size:         fileModel.Size,
		MimeType:     mimeType,
		Ext:          strings.TrimPrefix(ext, "."),
	}, nil
}
