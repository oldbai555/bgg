// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package daily_short_sentence

import (
	"context"
	"database/sql"
	"time"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyShortSentenceUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyShortSentenceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyShortSentenceUpdateLogic {
	return &DailyShortSentenceUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyShortSentenceUpdateLogic) DailyShortSentenceUpdate(req *types.DailyShortSentenceUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	sentence, err := l.svcCtx.Domain.Misc.DailyShortSentence.FindByID(l.ctx, req.Id)
	if err != nil {
		return errs.Wrap(errs.CodeNotFound, "每日短句不存在", err)
	}

	// 更新字段（只更新提供的字段）
	if req.SentenceType > 0 {
		sentence.Type = req.SentenceType
	}
	if req.Content != "" {
		sentence.Content = req.Content
	}
	if req.Img != "" {
		sentence.Img = sql.NullString{String: req.Img, Valid: true}
	} else if req.Img == "" && sentence.Img.Valid {
		// 如果传入空字符串，表示清空
		sentence.Img = sql.NullString{Valid: false}
	}
	if req.LiteratureAuthor != "" {
		sentence.LiteratureAuthor = sql.NullString{String: req.LiteratureAuthor, Valid: true}
	} else if req.LiteratureAuthor == "" && sentence.LiteratureAuthor.Valid {
		sentence.LiteratureAuthor = sql.NullString{Valid: false}
	}
	if req.ConvertImg != "" {
		sentence.ConvertImg = sql.NullString{String: req.ConvertImg, Valid: true}
	} else if req.ConvertImg == "" && sentence.ConvertImg.Valid {
		sentence.ConvertImg = sql.NullString{Valid: false}
	}

	sentence.UpdatedAt = time.Now().Unix()

	if err := l.svcCtx.Domain.Misc.DailyShortSentence.Update(l.ctx, sentence); err != nil {
		return errs.Wrap(errs.CodeInternalError, "更新每日短句失败", err)
	}

	return nil
}
