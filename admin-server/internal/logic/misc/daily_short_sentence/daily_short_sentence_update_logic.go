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

	_, err := l.svcCtx.IamRPC.DailyShortSentenceUpdate(l.ctx, &iamclient.DailyShortSentenceUpdateRequest{
		Id:               req.Id,
		SentenceType:     req.SentenceType,
		Content:          req.Content,
		Img:              req.Img,
		LiteratureAuthor: req.LiteratureAuthor,
		ConvertImg:       req.ConvertImg,
	})
	if err != nil {
		return errs.WrapGRPCError("更新每日短句失败", err)
	}
	return nil
}
