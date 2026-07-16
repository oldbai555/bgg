package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	videomodel "postapocgame/admin-server/services/content/internal/model/video"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoCollectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCollectLogic {
	return &VideoCollectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// VideoCollect 迁移自 internal/logic/video/video_collect/video_collect_logic.go。
// 原响应类型 types.VideoCollectResp 还带 Code/Msg 两个字段——和其余写接口一样，这两个
// 字段只是固定的成功文案，gateway 侧薄胶水自己拼装，不需要跨 RPC 边界传递。
func (l *VideoCollectLogic) VideoCollect(in *content.VideoCollectRequest) (*content.VideoCollectResponse, error) {
	if in.Uuid == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "新增失败: uuid 不能为空"))
	}

	existingVideo, err := l.svcCtx.Video.FindByUuid(l.ctx, in.Uuid)
	if err == nil && existingVideo != nil {
		return nil, toGRPCStatus(errs.New(errs.CodeConflict, "新增失败: 该 uuid 已存在"))
	}

	var xlzzUrlsJSON sql.NullString
	if len(in.XlzzUrls) > 0 {
		xlzzUrlsBytes, err := json.Marshal(in.XlzzUrls)
		if err != nil {
			return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "序列化磁力链接失败", err))
		}
		xlzzUrlsJSON = sql.NullString{String: string(xlzzUrlsBytes), Valid: true}
	}

	now := time.Now().Unix()
	video := &videomodel.Video{
		Uuid:      sql.NullString{String: in.Uuid, Valid: true},
		Name:      in.Name,
		PlayUrl:   in.PlayerUrl,
		Duration:  0,
		Type:      2, // 采集视频
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: 0,
	}
	if in.GodNum != "" {
		video.GodNum = sql.NullString{String: in.GodNum, Valid: true}
	}
	video.XlzzUrls = xlzzUrlsJSON

	if err := l.svcCtx.Video.Create(l.ctx, video); err != nil {
		l.Errorf("数据库插入失败: %v", err)
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "数据库插入失败", err))
	}

	return &content.VideoCollectResponse{Data: videoToItem(video)}, nil
}

func videoToItem(v *videomodel.Video) *content.VideoItem {
	item := &content.VideoItem{
		Id:         v.Id,
		Name:       v.Name,
		PlayUrl:    v.PlayUrl,
		Duration:   v.Duration,
		SourceType: v.Type,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
	}
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
