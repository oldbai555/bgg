package logic

import (
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeUpdateLogic {
	return &NoticeUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NoticeUpdateLogic) NoticeUpdate(in *iam.NoticeUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	notice, err := l.svcCtx.Domain.System.Notice.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "公告不存在", err))
	}

	if in.Title != "" {
		notice.Title = in.Title
	}
	if in.Content != "" {
		notice.Content = in.Content
	}
	if in.NoticeType > 0 {
		notice.Type = in.NoticeType
	}
	if in.Status > 0 {
		notice.Status = in.Status
	}
	if in.PublishTime > 0 {
		notice.PublishTime = in.PublishTime
	}

	oldStatus := notice.Status
	notice.UpdatedAt = time.Now().Unix()

	if err := l.svcCtx.Domain.System.Notice.Update(l.ctx, notice); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新公告失败", err))
	}

	if oldStatus == 1 && notice.Status == 2 {
		go l.createNotificationsForAllUsers(notice.Id, notice.Title, notice.Content)
	}

	return &iam.Empty{}, nil
}

func (l *NoticeUpdateLogic) createNotificationsForAllUsers(noticeID uint64, title, content string) {
	defer func() {
		if r := recover(); r != nil {
			l.Errorf("创建公告通知时发生 panic: %v, noticeId=%d", r, noticeID)
		}
	}()

	limit := int64(100)
	lastID := uint64(0)
	totalCreated := 0

	for {
		users, newLastID, err := l.svcCtx.Domain.IAM.User.FindChunk(context.Background(), limit, lastID)
		if err != nil {
			l.Errorf("查询用户失败: noticeId=%d, error: %v", noticeID, err)
			break
		}
		if len(users) == 0 {
			break
		}

		now := time.Now().Unix()
		for _, user := range users {
			notifications, _, err := l.svcCtx.Domain.System.Notification.FindPage(context.Background(), 1, 1, user.Id, "notice", -1)
			if err == nil {
				hasNotification := false
				for _, notif := range notifications {
					if notif.SourceId == noticeID && notif.SourceType == "notice" && notif.DeletedAt == 0 {
						hasNotification = true
						break
					}
				}
				if hasNotification {
					continue
				}
			}

			notification := &systemmodel.AdminNotification{
				UserId:     user.Id,
				SourceType: "notice",
				SourceId:   noticeID,
				Title:      title,
				Content:    content,
				ReadStatus: 1,
				ReadAt:     0,
				CreatedAt:  now,
				UpdatedAt:  now,
				DeletedAt:  0,
			}

			if err := l.svcCtx.Domain.System.Notification.Create(context.Background(), notification); err != nil {
				l.Errorf("创建公告通知失败: userId=%d, noticeId=%d, error: %v", user.Id, noticeID, err)
			} else {
				totalCreated++
			}
		}

		if len(users) < int(limit) {
			break
		}
		lastID = newLastID
	}

	l.Infof("成功为公告创建通知: noticeId=%d, totalCreated=%d", noticeID, totalCreated)
}
