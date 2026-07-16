package logic

import (
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/consts"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeCreateLogic {
	return &NoticeCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Notice / Notification
func (l *NoticeCreateLogic) NoticeCreate(in *iam.NoticeCreateRequest) (*iam.Empty, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}
	if in.Title == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "公告标题不能为空"))
	}
	if in.Content == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "公告内容不能为空"))
	}

	noticeType := in.NoticeType
	if noticeType == 0 {
		noticeType = 1
	}
	status := in.Status
	if status == 0 {
		status = consts.NoticeStatusDraft
	}
	publishTime := in.PublishTime
	if publishTime == 0 {
		publishTime = time.Now().Unix()
	}

	notice := &systemmodel.AdminNotice{
		Title:       in.Title,
		Content:     in.Content,
		Type:        noticeType,
		Status:      status,
		PublishTime: publishTime,
		CreatedBy:   in.OperatorUserId,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		DeletedAt:   0,
	}

	if err := l.svcCtx.Domain.System.Notice.Create(l.ctx, notice); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建公告失败", err))
	}

	if status == consts.NoticeStatusPublished {
		go l.createNotificationsForAllUsers(notice.Id, notice.Title, notice.Content)
	}

	return &iam.Empty{}, nil
}

// createNotificationsForAllUsers 给所有用户创建公告通知（User 和 Notification 现在同属
// iam-rpc 一个进程，不再需要跨服务回调）
func (l *NoticeCreateLogic) createNotificationsForAllUsers(noticeID uint64, title, content string) {
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
