// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package daily_short_sentence

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
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

	_, err := l.svcCtx.IamRPC.DailyShortSentenceCreate(l.ctx, &iamclient.DailyShortSentenceCreateRequest{
		SentenceType:     req.SentenceType,
		Content:          req.Content,
		Img:              req.Img,
		LiteratureAuthor: req.LiteratureAuthor,
		ConvertImg:       req.ConvertImg,
	})
	if err != nil {
		return errs.WrapGRPCError("创建每日短句失败", err)
	}
	return nil
}
