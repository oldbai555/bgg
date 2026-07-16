package logic

import (
	"context"
	"database/sql"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	miscmodel "postapocgame/admin-server/services/iam/internal/model/misc"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyShortSentenceCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDailyShortSentenceCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyShortSentenceCreateLogic {
	return &DailyShortSentenceCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DailyShortSentenceCreateLogic) DailyShortSentenceCreate(in *iam.DailyShortSentenceCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Content == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "短句内容不能为空"))
	}

	sentenceType := in.SentenceType
	if sentenceType == 0 {
		sentenceType = 1
	}

	now := time.Now().Unix()
	sentence := miscmodel.DailyShortSentence{
		Type:      sentenceType,
		Content:   in.Content,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: 0,
	}

	if in.Img != "" {
		sentence.Img = sql.NullString{String: in.Img, Valid: true}
	}
	if in.LiteratureAuthor != "" {
		sentence.LiteratureAuthor = sql.NullString{String: in.LiteratureAuthor, Valid: true}
	}
	if in.ConvertImg != "" {
		sentence.ConvertImg = sql.NullString{String: in.ConvertImg, Valid: true}
	}

	if err := l.svcCtx.Domain.Misc.DailyShortSentence.Create(l.ctx, &sentence); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建每日短句失败", err))
	}

	return &iam.Empty{}, nil
}
