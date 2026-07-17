#!/usr/bin/env bash
# =============================================================================
# Spec: 从 origin/main 切 Conventional Commits 类型分支、推送并用 gh 创建 PR
#
# Objective
#   自动化：拉取远端、必要时暂存本地改动、基于 main 最新创建命名分支、
#   push 到 origin，并向 main 发起 PR。
#
# 依赖
#   - git
#   - gh（GitHub CLI），需 gh auth login
#
# 命令示例
#   ./script/branch_from_main.sh
#   ./script/branch_from_main.sh --branch-type feat --slug order-export
#   ./script/branch_feat.sh order-export
#   ./script/branch_fix.sh order-export
#   ./script/branch_refactor.sh order-export
#
# 测试 / 验证
#   ./script/branch_from_main.sh --help
#   shellcheck script/branch_from_main.sh   # 若已安装
#
# Boundaries
#   - 不自动 git stash pop；不自动 commit。
#   - git fetch --all 会访问所有 remote，无权限的 remote 可能导致失败。
#   - 与 main 无领先提交时 PR 使用 --draft；仍失败则跳过 PR 并以 0 退出。
#   - squash 与否由仓库 PR 合并时的选择决定，本脚本不指定合并策略。
# =============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

REMOTE_NAME="origin"
TARGET_PR_BRANCH="main"
BRANCH_TYPE_PRESET=""
SLUG_RAW_PRESET=""

usage() {
	cat <<'EOF'
从 origin/main 创建分支、推送并创建 GitHub PR。

流程:
  git fetch origin →（若工作区有改动）git stash push -u → git fetch --all
  → 交互选择类型与描述（或使用 --branch-type + --slug）→ 新建分支并 push → gh pr create

分支类型（与 Conventional Commits 对齐）:
  feat      新功能
  fix       修复 bug
  docs      文档更新
  style     代码格式调整
  refactor  代码重构
  perf      性能优化
  test      测试相关
  chore     构建过程或辅助工具的变动

分支命名:
  <类型>/MMDD_<描述>（描述经清洗，仅字母数字与 _-）
  兼容别名：feature→feat、hotfix→fix

用法:
  ./script/branch_from_main.sh [选项]
  ./script/branch_feat.sh <slug> [选项]
  ./script/branch_fix.sh <slug> [选项]
  ./script/branch_docs.sh <slug> [选项]
  ./script/branch_style.sh <slug> [选项]
  ./script/branch_refactor.sh <slug> [选项]
  ./script/branch_perf.sh <slug> [选项]
  ./script/branch_test.sh <slug> [选项]
  ./script/branch_chore.sh <slug> [选项]

选项:
  --remote NAME              远端名（默认 origin）
  --branch-type TYPE         非交互：feat|fix|docs|style|refactor|perf|test|chore（或 feature|hotfix 别名）
  --slug TEXT                非交互：短描述，如 order-export（须与 --branch-type 同时使用）
  -h, --help                 显示本说明

说明:
  若曾自动 stash，请自行 git stash list / git stash pop 恢复。
  gh 命令若被 shell alias/函数劫持（例如个人 dotfiles 里给其它项目定义的同名别名），
  本脚本会在开头检测并报错提醒，避免误跑到别的脚本上。
EOF
}

die() {
	echo "错误: $*" >&2
	exit 1
}

require_cmd() {
	command -v "$1" >/dev/null 2>&1 || die "未找到命令: $1"
}

sanitize_slug() {
	local out
	out="$(printf '%s' "$1" | tr -cd 'a-zA-Z0-9_-')"
	printf '%s' "$out"
}

# 输入类型（含 feature/hotfix 别名）→ 分支前缀
normalize_branch_type() {
	case "$1" in
	feat | feature) printf '%s' 'feat' ;;
	fix | hotfix) printf '%s' 'fix' ;;
	docs | style | refactor | perf | test | chore) printf '%s' "$1" ;;
	*)
		return 1
		;;
	esac
}

valid_branch_types_help() {
	printf '%s' 'feat|fix|docs|style|refactor|perf|test|chore（feature|hotfix 为别名）'
}

while [[ $# -gt 0 ]]; do
	case "$1" in
	-h | --help)
		usage
		exit 0
		;;
	--remote)
		REMOTE_NAME="${2:?--remote 需要参数}"
		shift 2
		;;
	--branch-type)
		BRANCH_TYPE_PRESET="${2:?--branch-type 需要参数}"
		shift 2
		;;
	--slug)
		SLUG_RAW_PRESET="${2:?--slug 需要参数}"
		shift 2
		;;
	*)
		die "未知参数: $1（使用 --help）"
		;;
	esac
done

if [[ -n "$BRANCH_TYPE_PRESET" || -n "$SLUG_RAW_PRESET" ]]; then
	[[ -n "$BRANCH_TYPE_PRESET" && -n "$SLUG_RAW_PRESET" ]] ||
		die "--branch-type 与 --slug 必须同时指定（非交互模式）"
	BRANCH_TYPE="$(normalize_branch_type "$BRANCH_TYPE_PRESET")" ||
		die "无效的 --branch-type: ${BRANCH_TYPE_PRESET}（$(valid_branch_types_help)）"
fi

cd "$REPO_ROOT"

[[ -f AGENTS.md ]] || die "请在 bgg 仓库根目录执行（缺少 AGENTS.md）"
require_cmd git
git rev-parse --is-inside-work-tree >/dev/null 2>&1 || die "当前目录不是 git 工作区"

require_cmd gh
gh --version >/dev/null 2>&1 || die "gh 命令异常（可能被 shell alias/函数劫持成了别的脚本），请检查: type gh / which gh"
if ! gh auth status >/dev/null 2>&1; then
	die "gh 未登录或鉴权失败。请执行: gh auth login"
fi

BASE_REF="${REMOTE_NAME}/${TARGET_PR_BRANCH}"

echo "==> git fetch ${REMOTE_NAME}"
git fetch "$REMOTE_NAME"

if [[ -n "$(git status --porcelain)" ]]; then
	echo "==> 工作区有未提交改动，执行 git stash push -u"
	git stash push -u -m "branch_from_main: auto stash $(date +%Y%m%d-%H%M%S)"
	echo "提示: 结束后可用 git stash list / git stash pop 恢复"
fi

echo "==> git fetch --all"
git fetch --all

git rev-parse --verify "${BASE_REF}" >/dev/null 2>&1 || die "找不到基线引用: ${BASE_REF}（请先 fetch 并确认远端存在 main）"

if [[ -n "$BRANCH_TYPE_PRESET" ]]; then
	desc_raw="$SLUG_RAW_PRESET"
	echo "==> 非交互模式: ${BRANCH_TYPE} / ${desc_raw}"
else
	echo "选择分支类型:"
	echo "  1) feat      （新功能）"
	echo "  2) fix       （修复 bug）"
	echo "  3) docs      （文档更新）"
	echo "  4) style     （代码格式调整）"
	echo "  5) refactor  （代码重构）"
	echo "  6) perf      （性能优化）"
	echo "  7) test      （测试相关）"
	echo "  8) chore     （构建过程或辅助工具的变动）"
	read -r -p "请输入 1-8: " type_choice
	case "$type_choice" in
	1) BRANCH_TYPE="feat" ;;
	2) BRANCH_TYPE="fix" ;;
	3) BRANCH_TYPE="docs" ;;
	4) BRANCH_TYPE="style" ;;
	5) BRANCH_TYPE="refactor" ;;
	6) BRANCH_TYPE="perf" ;;
	7) BRANCH_TYPE="test" ;;
	8) BRANCH_TYPE="chore" ;;
	*)
		die "无效选择: ${type_choice}（请输入 1-8）"
		;;
	esac

	read -r -p "请输入短描述（如 order-export）: " desc_raw
fi
[[ -n "${desc_raw}" ]] || die "描述不能为空"

SLUG="$(sanitize_slug "$desc_raw")"
[[ -n "$SLUG" ]] || die "描述清洗后为空，请使用字母、数字、下划线或连字符"

DATE_MMdd="$(date +%m%d)"
BRANCH="${BRANCH_TYPE}/${DATE_MMdd}_${SLUG}"

if git show-ref --verify --quiet "refs/heads/${BRANCH}"; then
	die "本地已存在分支: ${BRANCH}"
fi
if git show-ref --verify --quiet "refs/remotes/${REMOTE_NAME}/${BRANCH}"; then
	die "远端已存在分支: ${REMOTE_NAME}/${BRANCH}"
fi

echo "==> 创建并切换到 ${BRANCH}（基于 ${BASE_REF}）"
git switch -c "${BRANCH}" "${BASE_REF}"

echo "==> git push -u ${REMOTE_NAME} HEAD"
git push -u "$REMOTE_NAME" HEAD

COMMITS_AHEAD="$(git rev-list --count "${BASE_REF}..HEAD" 2>/dev/null || echo 0)"
# 防御非数字
[[ "$COMMITS_AHEAD" =~ ^[0-9]+$ ]] || COMMITS_AHEAD=0

PR_TITLE="${BRANCH}"
PR_BODY="Created by script/branch_from_main.sh from ${BASE_REF}."

run_pr_create() {
	local use_draft="$1"
	local -a args=(
		pr create
		--base "$TARGET_PR_BRANCH"
		--head "$BRANCH"
		--title "$PR_TITLE"
		--body "$PR_BODY"
	)
	if [[ "$use_draft" -eq 1 ]]; then
		args+=(--draft)
	fi
	gh "${args[@]}"
}

echo "==> 创建 PR（领先提交数: ${COMMITS_AHEAD}）"

if [[ "$COMMITS_AHEAD" -eq 0 ]]; then
	echo "（无领先提交，将使用 --draft）"
	if run_pr_create 1; then
		echo "完成: PR 已创建（draft）。"
		exit 0
	fi
	echo "警告: gh pr create（draft）失败，可能为无 diff / 权限 / 项目策略。分支已推送到 ${REMOTE_NAME}。" >&2
	echo "可稍后手动执行:" >&2
	printf '  gh pr create --base %q --head %q --title %q --body %q --draft\n' \
		"$TARGET_PR_BRANCH" "$BRANCH" "$PR_TITLE" "$PR_BODY" >&2
	exit 0
fi

if run_pr_create 0; then
	echo "完成: PR 已创建。"
	exit 0
fi

echo "提示: 首次创建失败，尝试加 --draft 再试..." >&2
if run_pr_create 1; then
	echo "完成: PR 已创建（draft）。"
	exit 0
fi

echo "警告: gh pr create 仍失败。分支已推送到 ${REMOTE_NAME}。" >&2
echo "可稍后手动执行:" >&2
printf '  gh pr create --base %q --head %q --title %q --body %q\n' \
	"$TARGET_PR_BRANCH" "$BRANCH" "$PR_TITLE" "$PR_BODY" >&2
exit 0
