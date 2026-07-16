# 12 · scripts/ 规范化

## 前置依赖

无（这部分可以独立于 Phase 1 主线随时执行，属于低风险的规范化清理）。

## 0. 已核查的具体问题

以下每一条都已对照当前仓库实地核查（`admin-server/scripts/` 下 `generate-sql.sh`、`generate-model.sh`、`generate-api.sh`、`generate-ts.sh`、`migrate-menu.sh`、`README.md`、`sqlgen/main.go` 均已通读），不是照抄计划文档未经验证的判断。

### 0.1 `scripts/sqlgen/sqlgen.exe` 是被 git 跟踪的 Windows 编译产物

**核查结果**：`git ls-files admin-server/scripts/sqlgen/sqlgen.exe` 返回 `admin-server/scripts/sqlgen/sqlgen.exe`，确认**已被 git 跟踪提交**（文件本身 4.6MB，`ls -la` 确认）。`admin-server/.gitignore` 现有内容只排除了 `admin-server/*.exe`（仓库根一级），**不排除子目录 `scripts/sqlgen/` 下的 `.exe`**，所以每次 `generate-sql.sh` 在其他平台（非 Windows）重新 `go build -o sqlgen main.go` 产生同名产物时，`.gitignore` 也拦不住误提交（因为该规则只匹配 `sqlgen.exe`，不匹配 `sqlgen`，但 Windows 用户执行时会不小心把 `.exe` 后缀的历史遗留产物留在原地）。

**修复**：
1. `git rm --cached admin-server/scripts/sqlgen/sqlgen.exe`（只从 git 索引移除，不删本地文件，避免 Windows 开发者本地缺文件报错——实际上这个文件是可重新编译产物，删本地也没关系，但用 `--cached` 更保守）。
2. 在 `admin-server/.gitignore` 里新增一行 `scripts/sqlgen/sqlgen.exe`（或更通用的 `scripts/sqlgen/sqlgen` + `scripts/sqlgen/sqlgen.exe`，覆盖 `generate-sql.sh` 第 147 行 `go build -o sqlgen main.go` 在类 Unix 平台产生的无后缀产物和 Windows 平台产生的 `.exe` 产物两种情况）。当前 `.gitignore` 已有的 `admin-server/*.exe`/`**/*.exe` 这两条理论上能覆盖 `**/*.exe` 模式（`**/*.exe` 是全仓库通配），但 `sqlgen.exe` 依然被跟踪说明**它是在 `.gitignore` 加入 `**/*.exe` 规则之前就已经提交过的**——`.gitignore` 只对未跟踪文件生效，已跟踪文件即使匹配 `.gitignore` 规则也不会被自动忽略，这正是为什么需要显式 `git rm --cached` 这一步，光加 `.gitignore` 规则不够。

### 0.2 `scripts/README.md` 的示例与 `AGENTS.md` 当前规则矛盾

**核查结果**：`scripts/README.md` 第 40-46 行的示例：
```
./scripts/generate-sql.sh -group user -name 用户管理
./scripts/generate-sql.sh -group file -name 文件管理
./scripts/generate-sql.sh -group operation_log -name 操作日志 -parent-path /system
```
以及第 32、36 行的参数说明"`-group <group>`: 功能组名（必需，如 `user`, `file`）"，都是**扁平 group 名**，与 `AGENTS.md` 第 3 节现行规则"Group 格式 `<domain>/<module>`（如 `iam/user`、`blog/article`）"矛盾。`generate-sql.sh` 脚本本身（`main.go`/模板）**不限制** `-group` 的格式，纯字符串传递，技术上可以传 `iam/user`，只是 README 文档的示例没有跟上 DDD-lite 重构后的规则更新——这是文档滞后，不是脚本功能限制。

**修复**：改写 `scripts/README.md` 第 25-112 行涉及示例的部分：
- 第 36 行参数说明改为："`-group <group>`：功能组名（必需，格式 `<domain>/<module>`，snake_case，如 `iam/user`、`system/file`）"。
- 第 40-46 行示例改为：
  ```
  ./scripts/generate-sql.sh -group iam/user -name 用户管理
  ./scripts/generate-sql.sh -group system/file -name 文件管理
  ./scripts/generate-sql.sh -group monitoring/operation_log -name 操作日志 -parent-path /system
  ```
- 第 62-77 节"输出内容"里出现的 `init_<group>.sql`/`<group>.api.temp`/`<GroupUpper>List.vue` 等路径说明本身不需要改（沿用 `<group>` 占位符即可，实际值会是 `iam/user` 这种带斜杠的字符串，落到文件名/路径时脚本内部如何处理斜杠——需要核对 `sqlgen/main.go` 的路径拼接逻辑是否已经支持嵌套目录输出，如果没支持，这是脚本层面的 bug，需要一并修，不能只改文档不改代码让文档和实际行为继续对不上）。
- 第 90 行"注意事项"里如果有进一步引用旧示例格式的地方，一并核对更新。

### 0.3 `generate-sql.sh` 的退出码判断 bug

**核查结果**：`generate-sql.sh` 第 147 行 `go build -o sqlgen main.go` 编译，第 152 行 `./sqlgen ...` 实际运行程序，**第 155 行 `rm -f sqlgen` 清理编译产物**，**第 157 行 `if [ $? -eq 0 ]` 判断退出码**——但 `$?` 取的是**上一条命令**（第 155 行 `rm -f sqlgen`）的退出码，不是第 152 行 `sqlgen` 程序本身的退出码。`rm -f` 在文件不存在时也返回 0（`-f` 就是为了不因文件不存在报错），所以这个 `if` 分支**几乎恒为真**——即使第 152 行 `./sqlgen ...` 实际执行失败（比如参数错误、写文件权限问题），只要第 155 行的 `rm -f sqlgen` 成功执行，第 157-166 行依然会打印"✓ SQL 脚本生成成功!"这个具有误导性的成功提示。

**修复**：在第 152 行执行完 `./sqlgen ...` 后立刻保存其退出码到变量，再执行第 155 行的清理：
```bash
# 运行程序
./sqlgen -group "$GROUP" -name "$NAME" -output "$OUTPUT_DIR" -template "${SQLGEN_DIR}/templates" -parent-id "$PARENT_ID" -parent-path "$PARENT_PATH"
SQLGEN_EXIT_CODE=$?

# 清理编译产物
rm -f sqlgen

if [ $SQLGEN_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✓ SQL 脚本生成成功!${NC}"
    ...
else
    echo -e "${RED}✗ SQL 脚本生成失败${NC}"
    exit 1
fi
```
即：把第 152-166 行的"运行 → 清理 → 判断"顺序改成"运行 → 存退出码 → 清理 → 用存下来的退出码判断"。这是纯 bug 修复，不改变脚本对外行为契约（参数、输出文件路径都不变）。

## 1. 修复优先级与执行方式

三项都是低风险、范围明确的修复，按 `10-dev-execution-and-review-points.md` 的口径属于"本地开发相关的脚本/文档改动"，可以直接执行不用停下来问，修完作为正常 diff 走 review：

1. 先修 0.3（`generate-sql.sh` 的退出码 bug）——纯代码修复，最快验证（跑一次 `-group test/tmp -name 测试` 观察真实失败场景下是否正确报错）。
2. 再修 0.1（`sqlgen.exe` 从 git 移除 + `.gitignore` 补规则）——`git rm --cached` 操作本身要在这几项都改完、确认没有其他未提交改动混在一起时执行，避免一次提交里既有代码修复又有大文件删除，diff 不清晰。
3. 最后改 0.2（README 示例文本）——依赖 0.2 里提到的"`sqlgen/main.go` 是否已支持嵌套 group 路径输出"这一核查结论，如果发现不支持，这一项要连带一次小的代码改动（不只是改文档），实际执行时把这个依赖关系交代清楚，不要先改了文档、代码却没跟上。

## 2. 未来新增脚本：`generate-rpc.sh` 与 `generate-swagger.sh`

Phase 2（B.3）会新增 `generate-rpc.sh`（跑 `goctl rpc`，复用 `.template/rpc/*.tpl`），Phase 3（C.2）会新增 `generate-swagger.sh`（跑 `goctl api plugin -plugin goctl-swagger=...`，产出 `docs/openapi/admin-api.json`）。**这两个脚本的完整设计不在本篇范围内**，届时分别在 `16-rpc-conventions.md`（Phase 2）和 `20-api-docs-generation.md`（Phase 3）里详细展开。本篇只标记它们是本目录未来的"同类脚本"，落地时必须遵循与现有脚本一致的约定，不要另起一套风格：

- **`GOCTL_BIN` 解析模式**：参照 `generate-model.sh`/`generate-api.sh`/`generate-ts.sh` 已有的统一写法——优先读环境变量 `GOCTL_BIN`，其次 `command -v goctl`，再次尝试 `$(go env GOPATH)/bin/goctl`，都找不到则报错并给出安装提示。不要为新脚本发明另一套解析逻辑。
- **彩色输出**：复用现有的 `RED`/`GREEN`/`YELLOW`/`NC` 变量定义风格（四个脚本目前各自重复定义了一份相同的颜色变量，新脚本继续沿用同样的转义序列，保持终端输出观感一致；如果 Phase 1/2 顺手把这四个重复定义提取成一个 `scripts/_colors.sh` 公共文件被 `source`，新脚本也要跟着用公共文件，不要在提取之后还各写各的）。
- **用户亲自执行的政策**：这两个脚本和现有 `generate-*.sh` 一样，默认受 `AGENTS.md` 第 6 节"用户亲自执行"约束；但按 `10-dev-execution-and-review-points.md` 的口径，本轮重构开发期间可以由 AI 直接执行、事后 review，不需要额外重新论证一遍——这条策略天然覆盖新增的这两个脚本，不需要在它们的设计文档里重复声明一次例外。
- **前置条件检查与错误提示格式**：参照现有脚本"检查目录是否存在 → 检查前置工具是否安装 → 检查参数完整性 → 执行前打印配置摘要 → 执行"这套固定顺序和 `${RED}错误: ...${NC}` 的错误提示格式。

## 完成的定义

- `generate-sql.sh` 第 152-166 行的退出码判断修复，人工验证一次"故意传错误参数触发失败"和"正常参数触发成功"两种场景，确认提示信息与实际结果一致。
- `git ls-files admin-server/scripts/sqlgen/sqlgen.exe` 返回空（不再被跟踪），`.gitignore` 新增规则确认后续重新 `go build` 产生的同名文件不会被 `git status` 标记为待添加。
- `scripts/README.md` 里所有 `-group` 示例改为 `<domain>/<module>` 格式，且与 `sqlgen/main.go` 的实际路径处理行为一致（不是文档改了但代码没跟上导致文档描述的行为实际跑不通）。
