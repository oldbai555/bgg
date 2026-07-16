package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoDeleteLogic {
	return &VideoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// VideoDelete 迁移自 internal/logic/video/video/video_delete_logic.go。
func (l *VideoDeleteLogic) VideoDelete(in *content.VideoDeleteRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}
	if err := l.svcCtx.Video.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除视频失败", err))
	}
	return &content.Empty{}, nil
}
