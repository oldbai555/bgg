package fitness

import (
	"context"
	"strconv"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"

	"github.com/zeromicro/go-zero/core/logx"
)

// SuggestionView 待处理建议 + 模板名称。
type SuggestionView struct {
	Suggestion fitnessmodel.FitnessSuggestion
	FromName   string
	ToName     string
}

// CurrentSuggestion 返回本周待处理的建议。每周第一次打开手机端时按上周及之前的数据跑一次规则（懒生成），
// 等价于「每周一生成」，但不需要额外的定时任务，也不会给一整周没打开的人白算。
func (s *Service) CurrentSuggestion(ctx context.Context, userID uint64) (*SuggestionView, error) {
	today := s.Today()
	weekStart := WeekStart(today)
	if err := s.ensureWeeklySuggestion(ctx, userID, today, weekStart); err != nil {
		// 建议是锦上添花，算失败不影响页面主流程
		logx.WithContext(ctx).Errorf("fitness 生成每周建议失败: userId=%d, err=%v", userID, err)
	}
	sug, err := s.Suggestion.FindByUserWeek(ctx, userID, weekStart)
	if err != nil {
		return nil, dbErr(err)
	}
	if sug == nil || sug.Status != consts.FitnessSuggestionPending {
		return nil, nil
	}
	names, err := s.TemplateNames(ctx)
	if err != nil {
		return nil, err
	}
	return &SuggestionView{Suggestion: *sug, FromName: names[sug.FromTemplateId], ToName: names[sug.ToTemplateId]}, nil
}

func (s *Service) ensureWeeklySuggestion(ctx context.Context, userID uint64, today, weekStart string) error {
	key := consts.RedisFitnessSuggestCheckedPrefix + strconv.FormatUint(userID, 10) + ":" + weekStart
	if s.repo.Redis != nil {
		if ok, err := s.repo.Redis.ExistsCtx(ctx, key); err == nil && ok {
			return nil
		}
	}

	if err := s.expireOldSuggestions(ctx, userID, weekStart); err != nil {
		return err
	}
	existing, err := s.Suggestion.FindByUserWeek(ctx, userID, weekStart)
	if err != nil {
		return err
	}
	if existing == nil {
		input, current, ok, err := s.suggestInput(ctx, userID, today, weekStart)
		if err != nil {
			return err
		}
		if ok {
			if err := s.createSuggestion(ctx, userID, weekStart, current, Suggest(input)); err != nil {
				return err
			}
		}
	}

	if s.repo.Redis != nil {
		_ = s.repo.Redis.SetexCtx(ctx, key, "1", int((8 * 24 * time.Hour).Seconds()))
	}
	return nil
}

func (s *Service) expireOldSuggestions(ctx context.Context, userID uint64, weekStart string) error {
	pending, err := s.Suggestion.ListPendingByUser(ctx, userID)
	if err != nil {
		return err
	}
	for i := range pending {
		if pending[i].WeekStart < weekStart {
			pending[i].Status = consts.FitnessSuggestionExpired
			if err := s.Suggestion.Update(ctx, &pending[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

// suggestInput 组装规则输入。已经安排了次日切换的人不再给建议（ok=false）。
func (s *Service) suggestInput(ctx context.Context, userID uint64, today, weekStart string) (SuggestInput, fitnessmodel.FitnessTemplateSwitch, bool, error) {
	weeks := consts.FitnessUpgradeWeeks
	if consts.FitnessPlateauWeeks+1 > weeks {
		weeks = consts.FitnessPlateauWeeks + 1
	}
	start := AddDays(weekStart, -7*weeks)
	end := AddDays(weekStart, -1)

	facts, err := s.loadFacts(ctx, []uint64{userID}, start, end, true)
	if err != nil {
		return SuggestInput{}, fitnessmodel.FitnessTemplateSwitch{}, false, err
	}
	switches := facts.switches[userID]
	current, ok := TemplateAt(switches, today)
	if !ok {
		return SuggestInput{}, current, false, nil
	}
	if _, pending := PendingAfter(switches, today); pending {
		return SuggestInput{}, current, false, nil
	}
	tplMap, err := s.templateMap(ctx)
	if err != nil {
		return SuggestInput{}, current, false, err
	}
	profile, err := s.Profile.FindByUserID(ctx, userID)
	if err != nil {
		return SuggestInput{}, current, false, err
	}

	in := SuggestInput{
		CurrentCode:  tplMap[current.ToTemplateId].Code,
		CurrentSince: current.EffectiveDate,
		WeekStart:    weekStart,
	}
	if in.CurrentCode == consts.FitnessTemplatePlateau {
		in.PlateauReturnCode = tplMap[current.FromTemplateId].Code
	}
	if profile != nil {
		in.TargetWeight = profile.TargetWeightKg
	}

	points := make([]WeightPoint, 0)
	for _, b := range facts.bodies[userID] {
		points = append(points, WeightPoint{Date: b.RecordDate, Weight: b.WeightKg})
	}
	for k := weeks; k >= 1; k-- {
		ws := AddDays(weekStart, -7*k)
		we := AddDays(ws, 6)
		rate, counted := CompletionRate(facts.dayFacts(userID, DateRange(ws, we)), today)
		in.WeeklyRates = append(in.WeeklyRates, WeekRate{Rate: rate, Counted: counted})
		in.WeeklyAvgWeights = append(in.WeeklyAvgWeights, AvgWeightBetween(points, ws, we))
	}

	// 最新体重不限于最近几周：取最近一年里最后一次填写
	latest, err := s.BodyRecord.ListByUsersRange(ctx, []uint64{userID}, AddDays(today, -consts.FitnessBodyRecordMaxDays), today)
	if err != nil {
		return SuggestInput{}, current, false, err
	}
	for i := len(latest) - 1; i >= 0; i-- {
		if latest[i].WeightKg > 0 {
			in.LatestWeight = latest[i].WeightKg
			break
		}
	}
	return in, current, true, nil
}

func (s *Service) createSuggestion(ctx context.Context, userID uint64, weekStart string, current fitnessmodel.FitnessTemplateSwitch, hit *Suggestion) error {
	if hit == nil {
		return nil
	}
	to, err := s.Template.FindByCode(ctx, hit.ToCode)
	if err != nil {
		return err
	}
	if to == nil || to.Status != consts.FitnessStatusEnabled || to.Id == current.ToTemplateId {
		return nil
	}
	err = s.Suggestion.Create(ctx, &fitnessmodel.FitnessSuggestion{
		UserId:         userID,
		WeekStart:      weekStart,
		RuleCode:       hit.RuleCode,
		FromTemplateId: current.ToTemplateId,
		ToTemplateId:   to.Id,
		Reason:         truncate(hit.Reason, 255),
		Status:         consts.FitnessSuggestionPending,
	})
	if err != nil {
		// 并发生成撞 (user_id, week_start) 唯一键：另一个请求已经写好了
		if again, findErr := s.Suggestion.FindByUserWeek(ctx, userID, weekStart); findErr == nil && again != nil {
			return nil
		}
	}
	return err
}

// DecideSuggestion 用户对建议点「换」或「不换」。换 = 按建议发起一次次日生效的切换（来源：采纳建议）。
func (s *Service) DecideSuggestion(ctx context.Context, userID, id uint64, accept bool) error {
	sug, err := s.Suggestion.FindByID(ctx, id)
	if err != nil {
		return dbErr(err)
	}
	if sug == nil || sug.UserId != userID {
		return errs.New(errs.CodeNotFound, "建议不存在")
	}
	if sug.Status != consts.FitnessSuggestionPending {
		return errs.New(errs.CodeConflict, "这条建议已经处理过了")
	}
	if accept {
		if _, err := s.SwitchTemplate(ctx, userID, sug.ToTemplateId, consts.FitnessSwitchSourceSuggestion, consts.FitnessSwitchReasonSuggestionPrefix+sug.Reason, userID); err != nil {
			return err
		}
		sug.Status = consts.FitnessSuggestionAccepted
	} else {
		sug.Status = consts.FitnessSuggestionRejected
	}
	sug.DecidedAt = s.now().Unix()
	return dbErr(s.Suggestion.Update(ctx, sug))
}
