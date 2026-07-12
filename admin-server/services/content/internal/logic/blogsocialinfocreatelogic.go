package logic

import (
	"context"
	"fmt"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	blogmodel "postapocgame/admin-server/services/content/internal/model/blog"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogSocialInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogSocialInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogSocialInfoCreateLogic {
	return &BlogSocialInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogSocialInfoCreate 迁移自 internal/logic/blog/social_info/blog_social_info_create_logic.go。
func (l *BlogSocialInfoCreateLogic) BlogSocialInfoCreate(in *content.BlogSocialInfoCreateRequest) (*content.Empty, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "社交平台名称不能为空"))
	}
	url := strings.TrimSpace(in.Url)
	if url == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "目标链接不能为空"))
	}

	limits := l.svcCtx.Config.Limits
	if err := validateLength(name, limits.BlogSocialInfoNameMaxLength, "社交平台名称"); err != nil {
		return nil, toGRPCStatus(err)
	}
	if int64(len([]rune(url))) > limits.BlogSocialInfoUrlMaxLength {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, fmt.Sprintf("目标链接长度不能超过%d个字符", limits.BlogSocialInfoUrlMaxLength)))
	}

	remark := strings.TrimSpace(in.Remark)
	if remark != "" {
		if err := validateLength(remark, limits.BlogSocialInfoRemarkMaxLength, "备注"); err != nil {
			return nil, toGRPCStatus(err)
		}
	}

	if in.Status != 1 && in.Status != 2 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "状态值无效，必须为1（启用）或2（禁用）"))
	}

	info := &blogmodel.BlogSocialInfo{
		Name:     name,
		Url:      url,
		Remark:   remark,
		Status:   in.Status,
		OrderNum: in.OrderNum,
	}
	if err := l.svcCtx.SocialInfo.Create(l.ctx, info); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &content.Empty{}, nil
}
