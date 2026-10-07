package consts

// 身材管理（fitness 域）常量。数值型枚举与 db/services/iam/fitness/migrations/dict_fitness_20261007.sql
// 里的字典一一对应，字典 value 从 1 开始，0 表示「未填/不筛选」。

// 飞书登录场景（LoginFeishuRequest.scene）
const (
	FeishuSceneAdmin        = "admin"
	FeishuSceneFitnessOAuth = "fitness_oauth"
	FeishuSceneFitnessH5    = "fitness_h5"
)

// 模板/通用提示的启用状态（实体自身的状态列，不走字典）；0 只表示「未传」，新建时按启用处理
const (
	FitnessStatusEnabled  int64 = 1
	FitnessStatusDisabled int64 = 2
)

// 导出时把数值翻译成字典标签用到的字典编码
const (
	DictCodeFitnessTrainingType   = "fitness_training_type"
	DictCodeFitnessTrainingStatus = "fitness_training_status"
)

// 训练类型（字典 fitness_training_type）
const (
	FitnessTrainingStrength int64 = 1
	FitnessTrainingCardio   int64 = 2
	FitnessTrainingRest     int64 = 3
)

// 训练打卡状态（字典 fitness_training_status），0 表示当天还没打训练卡
const (
	FitnessTrainingStatusNone    int64 = 0
	FitnessTrainingStatusDone    int64 = 1
	FitnessTrainingStatusPartial int64 = 2
	FitnessTrainingStatusSkipped int64 = 3
)

// 三餐打卡状态（字典 fitness_meal_status）
const (
	FitnessMealStatusNone    int64 = 0
	FitnessMealStatusOnPlan  int64 = 1
	FitnessMealStatusOther   int64 = 2
	FitnessMealStatusSkipped int64 = 3
)

// 三餐时段编码，周模板 meals JSON 与打卡 meals JSON 共用
const (
	FitnessMealBreakfast = "breakfast"
	FitnessMealLunch     = "lunch"
	FitnessMealSnack     = "snack"
	FitnessMealDinner    = "dinner"
	FitnessMealAfterHome = "after_home"
)

// FitnessMealSlots 三餐时段的固定顺序
var FitnessMealSlots = []string{FitnessMealBreakfast, FitnessMealLunch, FitnessMealSnack, FitnessMealDinner, FitnessMealAfterHome}

// 训练时间偏好（字典 fitness_training_preference）
const (
	FitnessPreferenceMorning int64 = 1
	FitnessPreferenceEvening int64 = 2
)

// 模板切换来源（字典 fitness_switch_source）
const (
	FitnessSwitchSourceManual     int64 = 1
	FitnessSwitchSourceSuggestion int64 = 2
	FitnessSwitchSourceAdmin      int64 = 3
	FitnessSwitchSourceDefault    int64 = 4
)

// 切换记录的默认原因文案（管理员指定可以自己填原因，不填用默认）
const (
	FitnessSwitchReasonManual           = "手动切换"
	FitnessSwitchReasonAdmin            = "管理员指定"
	FitnessSwitchReasonSuggestionPrefix = "采纳建议："
)

// 模板建议状态（字典 fitness_suggestion_status）
const (
	FitnessSuggestionPending  int64 = 1
	FitnessSuggestionAccepted int64 = 2
	FitnessSuggestionRejected int64 = 3
	FitnessSuggestionExpired  int64 = 4
)

// 种子模板编码，切换建议规则按编码识别模板（管理员可以改名称/内容，但不要改这 7 个编码）
const (
	FitnessTemplateStarter  = "starter"
	FitnessTemplateStandard = "standard"
	FitnessTemplateAdvanced = "advanced"
	FitnessTemplateEvening  = "evening"
	FitnessTemplateBusy     = "busy"
	FitnessTemplatePlateau  = "plateau"
	FitnessTemplateMaintain = "maintain"
)

// FitnessDefaultTemplateCode 首次登录默认分配的模板：零基础入门。
// 新人第一周最大的风险是练太猛受伤/放弃，入门模板 3 次/周、自重为主，完成率高了再由建议规则升档。
const FitnessDefaultTemplateCode = FitnessTemplateStarter

// 建议规则编码（fitness_suggestion.rule_code）
const (
	FitnessRuleTargetReached = "target_reached"
	FitnessRulePlateauExit   = "plateau_exit"
	FitnessRuleDowngrade     = "downgrade"
	FitnessRuleUpgrade       = "upgrade"
	FitnessRulePlateau       = "plateau"
)

// 建议规则阈值（可调，改这里即可）
const (
	// FitnessUpgradeRate 连续 FitnessUpgradeWeeks 周训练完成率都 >= 该值时建议升一档
	FitnessUpgradeRate  = 0.8
	FitnessUpgradeWeeks = 4
	// FitnessDowngradeRate 连续 FitnessDowngradeWeeks 周训练完成率都 < 该值时建议降一档/换忙碌周
	FitnessDowngradeRate  = 0.5
	FitnessDowngradeWeeks = 2
	// FitnessPlateauWeeks 7 日均体重连续多少周不降时建议平台期模板
	FitnessPlateauWeeks = 2
	// FitnessPlateauStayWeeks 平台期模板用满多少周后建议回原模板
	FitnessPlateauStayWeeks = 2
)

// 打卡/统计口径
const (
	// FitnessBackfillDays 允许补卡的天数（今天往前 N 天）
	FitnessBackfillDays = 7
	// FitnessMovingAvgDays 体重移动平均窗口
	FitnessMovingAvgDays = 7
	// FitnessStreakLookbackDays 连续打卡天数最多往回查多少天
	FitnessStreakLookbackDays = 366
	// FitnessBodyRecordMaxDays 身体数据列表最多返回多少天
	FitnessBodyRecordMaxDays = 365
	// FitnessStatsMaxRangeDays 后台统计按日期筛选的最大跨度
	FitnessStatsMaxRangeDays = 366
	// FitnessRankingSize 后台连续打卡排行展示人数
	FitnessRankingSize = 10
)

// RedisFitnessSuggestCheckedPrefix 某人某周已经跑过建议规则的标记（无命中也要标记，避免每次打开页面都重算）
const RedisFitnessSuggestCheckedPrefix = "fitness:suggest:checked:"
