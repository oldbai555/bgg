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
	"postapocgame/admin-server/internal/model/misc"
	miscrepo "postapocgame/admin-server/internal/repository/misc"
)

type DailyShortSentenceCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyShortSentenceCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyShortSentenceCreateLogic {
	return &DailyShortSentenceCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyShortSentenceCreateLogic) DailyShortSentenceCreate(req *types.DailyShortSentenceCreateReq) error {
	if req == nil || req.Content == "" {
		return errs.New(errs.CodeBadRequest, "短句内容不能为空")
	}

	sentenceType := req.SentenceType
	if sentenceType == 0 {
		sentenceType = 1 // 默认为普通类型
	}

	now := time.Now().Unix()
	sentence := misc.DailyShortSentence{
		Type:      sentenceType,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: 0,
	}

	// 处理可选字段
	if req.Img != "" {
		sentence.Img = sql.NullString{String: req.Img, Valid: true}
	}
	if req.LiteratureAuthor != "" {
		sentence.LiteratureAuthor = sql.NullString{String: req.LiteratureAuthor, Valid: true}
	}
	if req.ConvertImg != "" {
		sentence.ConvertImg = sql.NullString{String: req.ConvertImg, Valid: true}
	}

	repo := miscrepo.NewDailyShortSentenceRepository(l.svcCtx.Repository)
	if err := repo.Create(l.ctx, &sentence); err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建每日短句失败", err)
	}

	return nil
}
