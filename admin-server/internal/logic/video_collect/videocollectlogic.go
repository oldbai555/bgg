// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video_collect

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
)

type VideoCollectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCollectLogic {
	return &VideoCollectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoCollectLogic) VideoCollect(req *types.VideoCollectReq) (resp *types.VideoCollectResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	// 验证 uuid 不为空
	if req.Uuid == "" {
		return nil, errs.New(errs.CodeBadRequest, "新增失败: uuid 不能为空")
	}

	// 验证 playerUrl 包含 uuid（安全检查，防止恶意请求）
	//if !strings.Contains(req.PlayerUrl, req.Uuid) {
	//	return nil, errs.New(errs.CodeBadRequest, "链接错误")
	//}

	// 检查 uuid 是否已存在
	videoRepo := repository.NewVideoRepository(l.svcCtx.Repository)
	existingVideo, err := videoRepo.FindByUuid(l.ctx, req.Uuid)
	if err == nil && existingVideo != nil {
		// uuid 已存在
		return nil, errs.New(errs.CodeConflict, "新增失败: 该 uuid 已存在")
	}

	// 序列化 xlzzUrls 为 JSON
	var xlzzUrlsJSON sql.NullString
	if len(req.XlzzUrls) > 0 {
		xlzzUrlsBytes, err := json.Marshal(req.XlzzUrls)
		if err != nil {
			return nil, errs.Wrap(errs.CodeInternalError, "序列化磁力链接失败", err)
		}
		xlzzUrlsJSON = sql.NullString{
			String: string(xlzzUrlsBytes),
			Valid:  true,
		}
	}

	// 构建 Video 模型
	now := time.Now().Unix()
	video := &model.Video{
		Uuid: sql.NullString{
			String: req.Uuid,
			Valid:  true,
		},
		Name:      req.Name,
		PlayUrl:   req.PlayerUrl,
		Duration:  0, // 默认值，采集时可能没有时长信息
		Type:      2, // 采集视频
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: 0,
	}

	// 处理可选字段
	if req.GodNum != "" {
		video.GodNum = sql.NullString{
			String: req.GodNum,
			Valid:  true,
		}
	}
	video.XlzzUrls = xlzzUrlsJSON

	// 插入数据
	err = videoRepo.Create(l.ctx, video)
	if err != nil {
		l.Errorf("数据库插入失败: %v", err)
		return nil, errs.Wrap(errs.CodeInternalError, "数据库插入失败", err)
	}

	// 转换为响应类型
	respItem := l.modelToVideoItem(video)

	return &types.VideoCollectResp{
		Code: 200,
		Msg:  "新增成功",
		Data: respItem,
	}, nil
}

// modelToVideoItem 将 model.Video 转换为 types.VideoItem
func (l *VideoCollectLogic) modelToVideoItem(v *model.Video) types.VideoItem {
	item := types.VideoItem{
		Id:         v.Id,
		Name:       v.Name,
		PlayUrl:    v.PlayUrl,
		Duration:   v.Duration,
		SourceType: v.Type,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
	}

	// 处理可选字段
	if v.Uuid.Valid {
		item.Uuid = v.Uuid.String
	}
	if v.Cover.Valid {
		item.Cover = v.Cover.String
	}
	if v.GodNum.Valid {
		item.GodNum = v.GodNum.String
	}
	if v.Description.Valid {
		item.Description = v.Description.String
	}
	if v.XlzzUrls.Valid {
		var xlzzUrls []string
		if err := json.Unmarshal([]byte(v.XlzzUrls.String), &xlzzUrls); err == nil {
			item.XlzzUrls = xlzzUrls
		}
	}

	return item
}
