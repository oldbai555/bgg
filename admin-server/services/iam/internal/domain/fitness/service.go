package fitness

import (
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	"postapocgame/admin-server/services/iam/internal/repository"
	fitnessrepo "postapocgame/admin-server/services/iam/internal/repository/fitness"

	"github.com/zeromicro/go-zero/core/logx"
)

// Service 身材管理领域服务：手机端与后台共用的业务规则都在这里，Logic 只做参数转换。
type Service struct {
	repo        *repository.Repository
	Template    fitnessrepo.TemplateRepository
	TemplateDay fitnessrepo.TemplateDayRepository
	Tip         fitnessrepo.TipRepository
	Profile     fitnessrepo.ProfileRepository
	Switch      fitnessrepo.SwitchRepository
	Checkin     fitnessrepo.CheckinRepository
	BodyRecord  fitnessrepo.BodyRecordRepository
	Suggestion  fitnessrepo.SuggestionRepository
	now         func() time.Time
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo:        repo,
		Template:    fitnessrepo.NewTemplateRepository(repo),
		TemplateDay: fitnessrepo.NewTemplateDayRepository(repo),
		Tip:         fitnessrepo.NewTipRepository(repo),
		Profile:     fitnessrepo.NewProfileRepository(repo),
		Switch:      fitnessrepo.NewSwitchRepository(repo),
		Checkin:     fitnessrepo.NewCheckinRepository(repo),
		BodyRecord:  fitnessrepo.NewBodyRecordRepository(repo),
		Suggestion:  fitnessrepo.NewSuggestionRepository(repo),
		now:         time.Now,
	}
}

// Today 上海时区的今天。
func (s *Service) Today() string {
	return Today(s.now())
}

func dbErr(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := errs.FromError(err); ok {
		return err
	}
	return errs.Wrap(errs.CodeBadDB, "数据库操作失败", err)
}

// EnsureMember 第一次使用身材管理时初始化档案 + 默认模板（当天生效）。幂等：登录时和每个手机端接口都会调用，
// 不依赖「必须先走过飞书登录」——后台账号直接打开手机端页面也能用。
func (s *Service) EnsureMember(ctx context.Context, userID uint64) (*fitnessmodel.FitnessProfile, error) {
	if userID == 0 {
		return nil, errs.New(errs.CodeUnauthorized, "未登录")
	}
	profile, err := s.Profile.FindByUserID(ctx, userID)
	if err != nil {
		return nil, dbErr(err)
	}
	if profile == nil {
		profile = &fitnessmodel.FitnessProfile{UserId: userID, TrainingPreference: consts.FitnessPreferenceMorning}
		if err := s.Profile.Create(ctx, profile); err != nil {
			// 并发首次请求撞唯一键：重新查一次自愈
			existing, findErr := s.Profile.FindByUserID(ctx, userID)
			if findErr != nil || existing == nil {
				return nil, dbErr(err)
			}
			profile = existing
		}
	}

	switches, err := s.Switch.ListByUsers(ctx, []uint64{userID})
	if err != nil {
		return nil, dbErr(err)
	}
	if len(switches) == 0 {
		tpl, err := s.Template.FindByCode(ctx, consts.FitnessDefaultTemplateCode)
		if err != nil {
			return nil, dbErr(err)
		}
		if tpl == nil {
			logx.WithContext(ctx).Errorf("fitness 默认模板 %s 不存在（未执行 init_fitness.sql？），userId=%d 暂无计划", consts.FitnessDefaultTemplateCode, userID)
			return profile, nil
		}
		if err := s.Switch.Create(ctx, &fitnessmodel.FitnessTemplateSwitch{
			UserId:        userID,
			ToTemplateId:  tpl.Id,
			EffectiveDate: EffectiveDateFor(consts.FitnessSwitchSourceDefault, s.Today()),
			Source:        consts.FitnessSwitchSourceDefault,
			Reason:        "首次使用，默认分配「" + tpl.Name + "」",
		}); err != nil {
			return nil, dbErr(err)
		}
	}
	return profile, nil
}

// SwitchTemplate 切换模板，次日生效（来源：手动/采纳建议/管理员指定）。
// 同一天内重复切换：先撤销所有尚未生效的切换，再按需写一条新的；切回当前模板等于撤销。
// 返回新模板的生效日期。
func (s *Service) SwitchTemplate(ctx context.Context, userID, templateID uint64, source int64, reason string, operatorID uint64) (string, error) {
	tpl, err := s.Template.FindByID(ctx, templateID)
	if err != nil {
		return "", dbErr(err)
	}
	if tpl == nil || tpl.Status != consts.FitnessStatusEnabled {
		return "", errs.New(errs.CodeBadRequest, "模板不存在或已停用")
	}
	today := s.Today()
	switches, err := s.Switch.ListByUsers(ctx, []uint64{userID})
	if err != nil {
		return "", dbErr(err)
	}
	current, hasCurrent := TemplateAt(switches, today)

	effective := EffectiveDateFor(source, today)
	err = s.repo.Transact(ctx, func(ctx context.Context, tx *repository.Repository) error {
		txSwitch := fitnessrepo.NewSwitchRepository(tx)
		for _, sw := range switches {
			if sw.EffectiveDate > today {
				if err := txSwitch.DeleteByID(ctx, sw.Id); err != nil {
					return err
				}
			}
		}
		if hasCurrent && current.ToTemplateId == templateID {
			effective = current.EffectiveDate
			return nil
		}
		return txSwitch.Create(ctx, &fitnessmodel.FitnessTemplateSwitch{
			UserId:         userID,
			FromTemplateId: current.ToTemplateId,
			ToTemplateId:   templateID,
			EffectiveDate:  effective,
			Source:         source,
			Reason:         truncate(reason, 255),
			OperatorId:     operatorID,
		})
	})
	if err != nil {
		return "", dbErr(err)
	}
	return effective, nil
}

// templateMap 全部模板按 ID 索引。
func (s *Service) templateMap(ctx context.Context) (map[uint64]fitnessmodel.FitnessTemplate, error) {
	list, err := s.Template.ListAll(ctx)
	if err != nil {
		return nil, dbErr(err)
	}
	out := make(map[uint64]fitnessmodel.FitnessTemplate, len(list))
	for _, t := range list {
		out[t.Id] = t
	}
	return out, nil
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func templateIDsOf(switches []fitnessmodel.FitnessTemplateSwitch) []uint64 {
	seen := map[uint64]bool{}
	var ids []uint64
	for _, sw := range switches {
		if sw.ToTemplateId > 0 && !seen[sw.ToTemplateId] {
			seen[sw.ToTemplateId] = true
			ids = append(ids, sw.ToTemplateId)
		}
	}
	return ids
}

func groupSwitches(list []fitnessmodel.FitnessTemplateSwitch) map[uint64][]fitnessmodel.FitnessTemplateSwitch {
	out := map[uint64][]fitnessmodel.FitnessTemplateSwitch{}
	for _, sw := range list {
		out[sw.UserId] = append(out[sw.UserId], sw)
	}
	return out
}
