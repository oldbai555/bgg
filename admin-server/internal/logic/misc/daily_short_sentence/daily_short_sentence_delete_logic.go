// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package daily_short_sentence

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	miscrepo "postapocgame/admin-server/internal/repository/misc"
)

type DailyShortSentenceDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyShortSentenceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyShortSentenceDeleteLogic {
	return &DailyShortSentenceDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyShortSentenceDeleteLogic) DailyShortSentenceDelete(req *types.DailyShortSentenceDeleteReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	repo := miscrepo.NewDailyShortSentenceRepository(l.svcCtx.Repository)
	if err := repo.DeleteByID(l.ctx, req.Id); err != nil {
		return errs.Wrap(errs.CodeInternalError, "删除每日短句失败", err)
	}

	return nil
}
