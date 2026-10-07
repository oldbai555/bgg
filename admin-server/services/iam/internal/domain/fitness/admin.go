package fitness

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	"postapocgame/admin-server/services/iam/internal/repository"
	fitnessrepo "postapocgame/admin-server/services/iam/internal/repository/fitness"
)

var codePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{1,31}$`)

func sortTemplates(list []fitnessmodel.FitnessTemplate) {
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].Sort != list[j].Sort {
			return list[i].Sort < list[j].Sort
		}
		return list[i].Id < list[j].Id
	})
}

// statusOrDefault 未传状态（0）时，新建按启用、编辑保持原值。
func statusOrDefault(status, fallback int64) int64 {
	if status == 0 {
		return fallback
	}
	return status
}

func validStatus(status int64) bool {
	return status == consts.FitnessStatusEnabled || status == consts.FitnessStatusDisabled
}

func validateTemplate(t *fitnessmodel.FitnessTemplate) error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" || len([]rune(t.Name)) > 64 {
		return errs.New(errs.CodeBadRequest, "模板名称不能为空且不超过 64 字")
	}
	if t.Level < 0 || t.Level > 9 {
		return errs.New(errs.CodeBadRequest, "层级需在 0~9 之间")
	}
	if !validStatus(t.Status) {
		return errs.New(errs.CodeBadRequest, "状态不合法")
	}
	if t.DailyKcal < 0 || t.DailyKcal > 6000 || t.DailyProtein < 0 || t.DailyProtein > 400 {
		return errs.New(errs.CodeBadRequest, "热量/蛋白质数值不合法")
	}
	t.Stage, t.Summary = truncate(strings.TrimSpace(t.Stage), 64), truncate(strings.TrimSpace(t.Summary), 500)
	return nil
}

func (s *Service) CreateTemplate(ctx context.Context, t *fitnessmodel.FitnessTemplate) error {
	t.Code = strings.TrimSpace(t.Code)
	if !codePattern.MatchString(t.Code) {
		return errs.New(errs.CodeBadRequest, "模板编码需为小写字母开头的字母/数字/下划线，2~32 位")
	}
	t.Status = statusOrDefault(t.Status, consts.FitnessStatusEnabled)
	if err := validateTemplate(t); err != nil {
		return err
	}
	existing, err := s.Template.FindByCode(ctx, t.Code)
	if err != nil {
		return dbErr(err)
	}
	if existing != nil {
		return errs.New(errs.CodeConflict, "模板编码已存在")
	}
	return dbErr(s.Template.Create(ctx, t))
}

// UpdateTemplate 编码创建后不可修改：建议规则按编码识别 7 套种子模板。
func (s *Service) UpdateTemplate(ctx context.Context, in *fitnessmodel.FitnessTemplate) error {
	existing, err := s.Template.FindByID(ctx, in.Id)
	if err != nil {
		return dbErr(err)
	}
	if existing == nil {
		return errs.New(errs.CodeNotFound, "模板不存在")
	}
	in.Code, in.CreatedAt = existing.Code, existing.CreatedAt
	in.Status = statusOrDefault(in.Status, existing.Status)
	if err := validateTemplate(in); err != nil {
		return err
	}
	if in.Status != consts.FitnessStatusEnabled {
		if err := s.ensureTemplateUnused(ctx, in.Id); err != nil {
			return err
		}
	}
	return dbErr(s.Template.Update(ctx, in))
}

// DeleteTemplate 软删模板及其每日内容；有人正在用（或明天起要用）的模板不能删/停用。
func (s *Service) DeleteTemplate(ctx context.Context, id uint64) error {
	existing, err := s.Template.FindByID(ctx, id)
	if err != nil {
		return dbErr(err)
	}
	if existing == nil {
		return errs.New(errs.CodeNotFound, "模板不存在")
	}
	if err := s.ensureTemplateUnused(ctx, id); err != nil {
		return err
	}
	rows, err := s.TemplateDay.ListByTemplates(ctx, []uint64{id}, "", "")
	if err != nil {
		return dbErr(err)
	}
	return dbErr(s.repo.Transact(ctx, func(ctx context.Context, tx *repository.Repository) error {
		txDay := fitnessrepo.NewTemplateDayRepository(tx)
		for _, row := range rows {
			if err := txDay.DeleteByID(ctx, row.Id); err != nil {
				return err
			}
		}
		return fitnessrepo.NewTemplateRepository(tx).DeleteByID(ctx, id)
	}))
}

func (s *Service) ensureTemplateUnused(ctx context.Context, templateID uint64) error {
	members, err := s.Profile.ListMembers(ctx, nil)
	if err != nil {
		return dbErr(err)
	}
	ids := make([]uint64, len(members))
	for i, m := range members {
		ids[i] = m.UserId
	}
	switches, err := s.Switch.ListByUsers(ctx, ids)
	if err != nil {
		return dbErr(err)
	}
	today := s.Today()
	inUse := 0
	for _, list := range groupSwitches(switches) {
		cur, ok := TemplateAt(list, today)
		pending, hasPending := PendingAfter(list, today)
		if ok && cur.ToTemplateId == templateID && !hasPending || hasPending && pending.ToTemplateId == templateID {
			inUse++
		}
	}
	if inUse > 0 {
		return errs.New(errs.CodeConflict, "还有使用者在用这个模板，请先在「使用者」里给他们换模板")
	}
	return nil
}

// TemplateDays 模板的周模板（固定 7 天，缺的补空）+ 全部日期覆盖。
func (s *Service) TemplateDays(ctx context.Context, templateID uint64) ([]fitnessmodel.FitnessTemplateDay, []fitnessmodel.FitnessTemplateDay, error) {
	rows, err := s.TemplateDay.ListByTemplates(ctx, []uint64{templateID}, "", "")
	if err != nil {
		return nil, nil, dbErr(err)
	}
	weekly := make([]fitnessmodel.FitnessTemplateDay, 7)
	for i := range weekly {
		weekly[i] = fitnessmodel.FitnessTemplateDay{TemplateId: templateID, Weekday: int64(i + 1), TrainingType: consts.FitnessTrainingRest, Exercises: "[]", Meals: "[]"}
	}
	var overrides []fitnessmodel.FitnessTemplateDay
	for _, row := range rows {
		if row.Weekday >= 1 && row.Weekday <= 7 {
			weekly[row.Weekday-1] = row
		} else if row.Weekday == 0 {
			overrides = append(overrides, row)
		}
	}
	return weekly, overrides, nil
}

// DayContent 编辑器提交的一天内容。
type DayContent struct {
	TrainingType int64
	TrainingTime string
	StepGoal     int64
	Warmup       string
	Exercises    []Exercise
	Stretch      string
	TrainingNote string
	Meals        []Meal
	Tip          string
}

// SaveTemplateDay 保存周模板某天（weekday 1~7）或某日期覆盖（weekday 0 + planDate），按业务键覆盖写。
func (s *Service) SaveTemplateDay(ctx context.Context, templateID uint64, weekday int64, planDate string, c DayContent) error {
	tpl, err := s.Template.FindByID(ctx, templateID)
	if err != nil {
		return dbErr(err)
	}
	if tpl == nil {
		return errs.New(errs.CodeNotFound, "模板不存在")
	}
	switch {
	case weekday >= 1 && weekday <= 7:
		planDate = ""
	case weekday == 0:
		if !ValidDate(planDate) {
			return errs.New(errs.CodeBadRequest, "按日期调整需要填写合法日期")
		}
	default:
		return errs.New(errs.CodeBadRequest, "星期需在 1~7 之间")
	}
	if c.TrainingType < consts.FitnessTrainingStrength || c.TrainingType > consts.FitnessTrainingRest {
		return errs.New(errs.CodeBadRequest, "训练类型不合法")
	}
	if c.StepGoal < 0 || c.StepGoal > 100000 {
		return errs.New(errs.CodeBadRequest, "步数目标不合法")
	}
	if len(c.Exercises) > 30 {
		return errs.New(errs.CodeBadRequest, "动作最多 30 个")
	}
	exercises := make([]Exercise, 0, len(c.Exercises))
	for _, e := range c.Exercises {
		e.Name = strings.TrimSpace(e.Name)
		if e.Name == "" {
			continue
		}
		if e.Sets < 0 || e.Sets > 20 {
			return errs.New(errs.CodeBadRequest, "组数需在 0~20 之间")
		}
		exercises = append(exercises, Exercise{Name: truncate(e.Name, 64), Sets: e.Sets, Reps: truncate(e.Reps, 32), Duration: truncate(e.Duration, 32), Note: truncate(e.Note, 200)})
	}
	meals, err := normalizeMeals(c.Meals)
	if err != nil {
		return err
	}

	row, err := s.TemplateDay.FindByKey(ctx, templateID, weekday, planDate)
	if err != nil {
		return dbErr(err)
	}
	isNew := row == nil
	if isNew {
		row = &fitnessmodel.FitnessTemplateDay{TemplateId: templateID, Weekday: weekday, PlanDate: planDate}
	}
	row.TrainingType = c.TrainingType
	row.TrainingTime = truncate(strings.TrimSpace(c.TrainingTime), 32)
	row.StepGoal = c.StepGoal
	row.Warmup = truncate(strings.TrimSpace(c.Warmup), 255)
	row.Exercises = EncodeJSON(exercises)
	row.Stretch = truncate(strings.TrimSpace(c.Stretch), 255)
	row.TrainingNote = truncate(strings.TrimSpace(c.TrainingNote), 500)
	row.Meals = EncodeJSON(meals)
	row.Tip = truncate(strings.TrimSpace(c.Tip), 1000)
	if isNew {
		return dbErr(s.TemplateDay.Create(ctx, row))
	}
	return dbErr(s.TemplateDay.Update(ctx, row))
}

// normalizeMeals 按固定餐次顺序输出，同一餐次只保留第一条。
func normalizeMeals(in []Meal) ([]Meal, error) {
	bySlot := map[string]Meal{}
	for _, m := range in {
		valid := false
		for _, slot := range consts.FitnessMealSlots {
			if m.Slot == slot {
				valid = true
				break
			}
		}
		if !valid {
			return nil, errs.New(errs.CodeBadRequest, "餐次不合法: "+m.Slot)
		}
		if _, ok := bySlot[m.Slot]; ok {
			continue
		}
		if m.Kcal < 0 || m.Kcal > 3000 || m.Protein < 0 || m.Protein > 300 {
			return nil, errs.New(errs.CodeBadRequest, "热量/蛋白质数值不合法")
		}
		alts := make([]string, 0, len(m.Alternatives))
		for _, a := range m.Alternatives {
			if a = strings.TrimSpace(a); a != "" && len(alts) < 5 {
				alts = append(alts, truncate(a, 100))
			}
		}
		bySlot[m.Slot] = Meal{
			Slot: m.Slot, Time: truncate(strings.TrimSpace(m.Time), 32), Place: truncate(strings.TrimSpace(m.Place), 32),
			Food: truncate(strings.TrimSpace(m.Food), 200), Portion: truncate(strings.TrimSpace(m.Portion), 100),
			Kcal: m.Kcal, Protein: m.Protein, HowToOrder: truncate(strings.TrimSpace(m.HowToOrder), 300), Alternatives: alts,
		}
	}
	out := make([]Meal, 0, len(bySlot))
	for _, slot := range consts.FitnessMealSlots {
		if m, ok := bySlot[slot]; ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// DeleteTemplateDay 只允许删日期覆盖；周模板 7 天是固定骨架，只能改不能删。
func (s *Service) DeleteTemplateDay(ctx context.Context, id uint64) error {
	row, err := s.TemplateDay.FindByID(ctx, id)
	if err != nil {
		return dbErr(err)
	}
	if row == nil {
		return errs.New(errs.CodeNotFound, "记录不存在")
	}
	if row.Weekday != 0 {
		return errs.New(errs.CodeBadRequest, "周模板只能修改不能删除")
	}
	return dbErr(s.TemplateDay.DeleteByID(ctx, id))
}

func validateTip(t *fitnessmodel.FitnessTip) error {
	t.Title, t.Content = strings.TrimSpace(t.Title), strings.TrimSpace(t.Content)
	if t.Title == "" || len([]rune(t.Title)) > 64 {
		return errs.New(errs.CodeBadRequest, "标题不能为空且不超过 64 字")
	}
	if t.Content == "" || len([]rune(t.Content)) > 2000 {
		return errs.New(errs.CodeBadRequest, "内容不能为空且不超过 2000 字")
	}
	if !validStatus(t.Status) {
		return errs.New(errs.CodeBadRequest, "状态不合法")
	}
	return nil
}

func (s *Service) CreateTip(ctx context.Context, t *fitnessmodel.FitnessTip) error {
	t.Code = strings.TrimSpace(t.Code)
	if !codePattern.MatchString(t.Code) {
		return errs.New(errs.CodeBadRequest, "编码需为小写字母开头的字母/数字/下划线，2~32 位")
	}
	t.Status = statusOrDefault(t.Status, consts.FitnessStatusEnabled)
	if err := validateTip(t); err != nil {
		return err
	}
	existing, err := s.Tip.FindByCode(ctx, t.Code)
	if err != nil {
		return dbErr(err)
	}
	if existing != nil {
		return errs.New(errs.CodeConflict, "编码已存在")
	}
	return dbErr(s.Tip.Create(ctx, t))
}

func (s *Service) UpdateTip(ctx context.Context, in *fitnessmodel.FitnessTip) error {
	existing, err := s.Tip.FindByID(ctx, in.Id)
	if err != nil {
		return dbErr(err)
	}
	if existing == nil {
		return errs.New(errs.CodeNotFound, "提示不存在")
	}
	in.Code, in.CreatedAt = existing.Code, existing.CreatedAt
	in.Status = statusOrDefault(in.Status, existing.Status)
	if err := validateTip(in); err != nil {
		return err
	}
	return dbErr(s.Tip.Update(ctx, in))
}

func (s *Service) DeleteTip(ctx context.Context, id uint64) error {
	existing, err := s.Tip.FindByID(ctx, id)
	if err != nil {
		return dbErr(err)
	}
	if existing == nil {
		return errs.New(errs.CodeNotFound, "提示不存在")
	}
	return dbErr(s.Tip.DeleteByID(ctx, id))
}

// AssignTemplate 管理员给某人指定模板（次日生效，来源：管理员指定）。
func (s *Service) AssignTemplate(ctx context.Context, userID, templateID uint64, reason string, operatorID uint64) (string, error) {
	member, err := s.Profile.FindByUserID(ctx, userID)
	if err != nil {
		return "", dbErr(err)
	}
	if member == nil {
		return "", errs.New(errs.CodeNotFound, "使用者不存在")
	}
	if strings.TrimSpace(reason) == "" {
		reason = consts.FitnessSwitchReasonAdmin
	}
	return s.SwitchTemplate(ctx, userID, templateID, consts.FitnessSwitchSourceAdmin, reason, operatorID)
}
