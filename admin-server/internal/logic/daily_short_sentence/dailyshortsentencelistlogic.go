// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package daily_short_sentence

import (
	"context"
	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

	repo := repository.NewDailyShortSentenceRepository(l.svcCtx.Repository)
	list, total, err := repo.FindPage(l.ctx, req.Page, req.PageSize, req.Keyword, req.SentenceType)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询每日短句列表失败", err)
	}

	items := make([]types.DailyShortSentenceItem, 0, len(list))
	for _, s := range list {
		item := types.DailyShortSentenceItem{
			Id:           s.Id,
			SentenceType: s.Type,
			Content:      s.Content,
			CreatedAt:    s.CreatedAt,
			UpdatedAt:    s.UpdatedAt,
		}

		// 处理可选字段
		if s.Img.Valid {
			item.Img = s.Img.String
		}
		if s.LiteratureAuthor.Valid {
			item.LiteratureAuthor = s.LiteratureAuthor.String
		}
		if s.ConvertImg.Valid {
			item.ConvertImg = s.ConvertImg.String
		}

		items = append(items, item)
	}

	return &types.DailyShortSentenceListResp{
		Total: total,
		List:  items,
	}, nil
}
