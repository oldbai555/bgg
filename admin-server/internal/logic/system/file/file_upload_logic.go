// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/model/system"
	systemrepo "postapocgame/admin-server/internal/repository/system"
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
	if err := os.MkdirAll(consts.UploadDir, 0755); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "创建上传目录失败", err)
	}

	// 计算文件的 MD5 哈希值
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "计算文件MD5失败", err)
	}
	md5Hash := fmt.Sprintf("%x", hash.Sum(nil))

	// 重置文件指针，以便后续读取
	file.Seek(0, 0)

	// 获取文件扩展名
	ext := filepath.Ext(header.Filename)
	// 使用 MD5 + 扩展名作为文件名
	fileName := md5Hash + ext

	// 获取基础 URL（从字典中读取）
	baseURL := l.getStorageBaseURL()

	// 检查文件是否已存在（根据 MD5）
	fileRepo := systemrepo.NewFileRepository(l.svcCtx.Repository)
	existingFile, err := fileRepo.FindByName(l.ctx, fileName)
	if err == nil && existingFile != nil {
		// 文件已存在，直接返回已有记录
		l.Infof("文件已存在，MD5: %s", md5Hash)
		proxyPath := fmt.Sprintf("%s/%s", consts.PathFileUploads, fileName)
		fullURL := proxyPath
		if baseURL != "" {
			if strings.HasPrefix(baseURL, "http://") || strings.HasPrefix(baseURL, "https://") {
				fullURL = fmt.Sprintf("%s%s", baseURL, proxyPath)
			}
		}

		return &types.FileUploadResp{
			Id:           existingFile.Id,
			Name:         existingFile.Name,
			OriginalName: existingFile.OriginalName,
			Path:         proxyPath,
			BaseUrl:      existingFile.BaseUrl,
			Url:          fullURL,
			Size:         existingFile.Size,
			MimeType:     existingFile.MimeType.String,
			Ext:          existingFile.Ext.String,
		}, nil
	}

	// 文件不存在，保存新文件
	fileSystemPath := filepath.Join(consts.UploadDir, fileName)
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
	fileModel := system.AdminFile{
		Name:         fileName, // MD5 + 扩展名
		OriginalName: header.Filename,
		Path:         fmt.Sprintf("%s/%s", consts.PathFileUploads, fileName), // 访问路径
		BaseUrl:      baseURL,                                                // 基础 URL
		Size:         uint64(fileInfo.Size()),
		MimeType:     sql.NullString{String: mimeType, Valid: mimeType != ""},
		Ext:          sql.NullString{String: strings.TrimPrefix(ext, "."), Valid: ext != ""},
		StorageType:  "local",
		Status:       1,
	}

	if err := fileRepo.Create(l.ctx, &fileModel); err != nil {
		// 如果数据库保存失败，删除已上传的文件
		os.Remove(fileSystemPath)
		return nil, errs.Wrap(errs.CodeInternalError, "保存文件记录失败", err)
	}

	// 返回文件访问路径
	proxyPath := fmt.Sprintf("%s/%s", consts.PathFileUploads, fileName)
	fullURL := proxyPath
	if baseURL != "" {
		if strings.HasPrefix(baseURL, "http://") || strings.HasPrefix(baseURL, "https://") {
			fullURL = fmt.Sprintf("%s%s", baseURL, proxyPath)
		}
	}

	return &types.FileUploadResp{
		Id:           fileModel.Id,
		Name:         fileModel.Name, // MD5 + 扩展名
		OriginalName: fileModel.OriginalName,
		Path:         proxyPath,
		BaseUrl:      baseURL,
		Url:          fullURL,
		Size:         fileModel.Size,
		MimeType:     mimeType,
		Ext:          strings.TrimPrefix(ext, "."),
	}, nil
}

// getStorageBaseURL 从字典中获取存储baseURL
func (l *FileUploadLogic) getStorageBaseURL() string {
	// 从字典中获取配置
	dictTypeRepo := systemrepo.NewDictTypeRepository(l.svcCtx.Repository)
	dictType, err := dictTypeRepo.FindByCode(l.ctx, "storage_base_url")
	if err != nil {
		l.Errorf("获取存储配置字典类型失败: %v", err)
		return ""
	}

	dictItemRepo := systemrepo.NewDictItemRepository(l.svcCtx.Repository)
	items, err := dictItemRepo.FindByTypeID(l.ctx, dictType.Id)
	if err != nil || len(items) == 0 {
		l.Errorf("获取存储配置字典项失败: %v", err)
		return ""
	}

	// 使用第一个有效的字典项值
	baseURL := items[0].Value
	if baseURL == "" {
		l.Errorf("字典中的baseURL为空")
		return ""
	}

	return strings.TrimSuffix(baseURL, "/")
}
