package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyShortSentenceListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDailyShortSentenceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyShortSentenceListLogic {
	return &DailyShortSentenceListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DailyShortSentenceListLogic) DailyShortSentenceList(in *iam.DailyShortSentenceListRequest) (*iam.DailyShortSentenceListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	page, pageSize := in.Page, in.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	list, total, err := l.svcCtx.Domain.Misc.DailyShortSentence.FindPage(l.ctx, page, pageSize, in.Keyword, in.SentenceType)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询每日短句列表失败", err))
	}

	items := make([]*iam.DailyShortSentenceItem, 0, len(list))
	for _, s := range list {
		item := &iam.DailyShortSentenceItem{
			Id:           s.Id,
			SentenceType: s.Type,
			Content:      s.Content,
			CreatedAt:    s.CreatedAt,
			UpdatedAt:    s.UpdatedAt,
		}
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

	return &iam.DailyShortSentenceListResponse{Total: total, List: items}, nil
}
