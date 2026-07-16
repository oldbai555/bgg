package logic

import (
	"context"
	"fmt"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogSocialInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogSocialInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogSocialInfoUpdateLogic {
	return &BlogSocialInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogSocialInfoUpdate 迁移自 internal/logic/blog/social_info/blog_social_info_update_logic.go。
func (l *BlogSocialInfoUpdateLogic) BlogSocialInfoUpdate(in *content.BlogSocialInfoUpdateRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "社交信息ID不能为空"))
	}

	info, err := l.svcCtx.SocialInfo.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询社交信息失败", err))
	}
	if info == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "社交信息不存在"))
	}

	limits := l.svcCtx.Config.Limits

	if in.Name != "" {
		name := strings.TrimSpace(in.Name)
		if err := validateLength(name, limits.BlogSocialInfoNameMaxLength, "社交平台名称"); err != nil {
			return nil, toGRPCStatus(err)
		}
		info.Name = name
	}
	if in.Url != "" {
		url := strings.TrimSpace(in.Url)
		if int64(len([]rune(url))) > limits.BlogSocialInfoUrlMaxLength {
			return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, fmt.Sprintf("目标链接长度不能超过%d个字符", limits.BlogSocialInfoUrlMaxLength)))
		}
		info.Url = url
	}
	if in.Remark != "" {
		remark := strings.TrimSpace(in.Remark)
		if err := validateLength(remark, limits.BlogSocialInfoRemarkMaxLength, "备注"); err != nil {
			return nil, toGRPCStatus(err)
		}
		info.Remark = remark
	}
	if in.Status > 0 {
		if in.Status != 1 && in.Status != 2 {
			return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "状态值无效，必须为1（启用）或2（禁用）"))
		}
		info.Status = in.Status
	}
	if in.OrderNum > 0 {
		info.OrderNum = in.OrderNum
	}

	if err := l.svcCtx.SocialInfo.Update(l.ctx, info); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &content.Empty{}, nil
}
