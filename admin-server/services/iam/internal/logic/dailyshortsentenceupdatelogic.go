package logic

import (
	"context"
	"database/sql"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyShortSentenceUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDailyShortSentenceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyShortSentenceUpdateLogic {
	return &DailyShortSentenceUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DailyShortSentenceUpdateLogic) DailyShortSentenceUpdate(in *iam.DailyShortSentenceUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	sentence, err := l.svcCtx.Domain.Misc.DailyShortSentence.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "每日短句不存在", err))
	}

	if in.SentenceType > 0 {
		sentence.Type = in.SentenceType
	}
	if in.Content != "" {
		sentence.Content = in.Content
	}
	if in.Img != "" {
		sentence.Img = sql.NullString{String: in.Img, Valid: true}
	} else if sentence.Img.Valid {
		sentence.Img = sql.NullString{Valid: false}
	}
	if in.LiteratureAuthor != "" {
		sentence.LiteratureAuthor = sql.NullString{String: in.LiteratureAuthor, Valid: true}
	} else if sentence.LiteratureAuthor.Valid {
		sentence.LiteratureAuthor = sql.NullString{Valid: false}
	}
	if in.ConvertImg != "" {
		sentence.ConvertImg = sql.NullString{String: in.ConvertImg, Valid: true}
	} else if sentence.ConvertImg.Valid {
		sentence.ConvertImg = sql.NullString{Valid: false}
	}

	sentence.UpdatedAt = time.Now().Unix()

	if err := l.svcCtx.Domain.Misc.DailyShortSentence.Update(l.ctx, sentence); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新每日短句失败", err))
	}

	return &iam.Empty{}, nil
}
