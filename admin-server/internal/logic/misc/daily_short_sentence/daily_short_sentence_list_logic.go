// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package daily_short_sentence

import (
	"context"

	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyShortSentenceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyShortSentenceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyShortSentenceListLogic {
	return &DailyShortSentenceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyShortSentenceListLogic) DailyShortSentenceList(req *types.DailyShortSentenceListReq) (resp *types.DailyShortSentenceListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)

	rpcResp, err := l.svcCtx.IamRPC.DailyShortSentenceList(l.ctx, &iamclient.DailyShortSentenceListRequest{
		Page:         req.Page,
		PageSize:     req.PageSize,
		Keyword:      req.Keyword,
		SentenceType: req.SentenceType,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询每日短句列表失败", err)
	}

	items := make([]types.DailyShortSentenceItem, 0, len(rpcResp.List))
	for _, s := range rpcResp.List {
		items = append(items, types.DailyShortSentenceItem{
			Id:               s.Id,
			SentenceType:     s.SentenceType,
			Content:          s.Content,
			Img:              s.Img,
			LiteratureAuthor: s.LiteratureAuthor,
			ConvertImg:       s.ConvertImg,
			CreatedAt:        s.CreatedAt,
			UpdatedAt:        s.UpdatedAt,
		})
	}

	return &types.DailyShortSentenceListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
