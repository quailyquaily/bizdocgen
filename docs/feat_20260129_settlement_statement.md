# 结算单（Settlement Statement）支持（替代 PaymentStatement）

## 背景
当前仓库已有 `paymentstatement` 体系，但需求是新增 **Settlement Statement**，并且 **完全沿用 invoice 的样式与模板**，仅文案替换（例如 Bill To -> Recipient、Invoice ID -> Statement ID 等）。同时决定 **完全删除** paymentstatement 支持，避免双轨维护。

## 目标
- 对外提供语义正确的接口：`SettlementStatement` 而非 `Invoice` / `PaymentStatement`。
- 结算单沿用 invoice 的布局模板（classic / modern / compact / spotlight / ledger / split）。
- 仅替换文案与标题；布局、结构与渲染逻辑尽量复用。
- 移除 paymentstatement 相关类型、builder、layout、样例、文档。

## 非目标
- 不再保留 paymentstatement 的任何兼容层或旧 API。
- 不引入复杂的跨文档抽象层（如大范围 `core/common` 接口）。
- 不改变 invoice 的字段结构或计算逻辑。

## 设计决策（评估要点）
- 需要语义清晰的对外接口，因此 **新增 SettlementStatement 对外入口**。
- 结算单与 invoice 仅文案差异，**无需新的一套字段结构**；避免过度抽象。
- 只提取极少量的共享辅助函数（如 label key 选择）来复用 invoice 的布局逻辑。

## API 与数据结构
- 新增：
  - `core/settlementstatement.go`（建议：类型别名或轻量 wrapper，字段与 `InvoiceParams` 保持一致）。
  - `builder/settlementstatement.go`：对外提供 `NewSettlementStatementBuilder` / `GenerateSettlementStatement`。
  - `builder.Config` 新增 `SettlementStatementLayout`（或复用 `InvoiceLayout`，但对外字段名称仍为 `SettlementStatementLayout`）。
- 删除：
  - `core/paymentstatement.go`
  - `builder/paymentstatement.go`
  - `builder/layout_paymentstatement.go`
  - `GeneratePaymentStatement` / `NewPaymentStatementBuilder` / `PaymentStatementLayout` 等相关入口

## 文案替换（i18n）
新增 settlement statement 对应的 label keys：
- `StatementID`
- `StatementIssueDate`
- `StatementPeriod`
- `StatementRecipient`
- （可选）`StatementTitle` / `StatementDescription`

并在所有语言文件中补齐：
- `i18n/locales/en.toml`
- `i18n/locales/ja.toml`
- `i18n/locales/zh_cn.toml`
- `i18n/locales/zh_tw.toml`

实现方式：
- builder 在生成 invoice 头部/收件人区域时，根据“文案集”选择 `Invoice*` 或 `Statement*`。
- 当文案集为 `statement` 时，键名以 `Statement` 开头；找不到则回退 `Invoice`。

## Builder 复用策略
- 直接复用 invoice 的 layout（classic / modern / compact / spotlight / ledger / split）。
- 在渲染处切换 label key（比如 ID、Issue Date、Period、Bill To）。
- 禁止引入大范围 `core/common` 或 `builder/common` 接口抽象；仅允许小型 helper。

## 样例与文档
- 删除 `samples/paymentstatement-*.yaml` 与生成逻辑。
- 新增 `samples/settlementstatement-1.yaml`（结构与 invoice 一致，文案集标记为 `statement`）。
- 更新 `cmd/generate-samples`：生成 settlement statement 示例 PDF。
- 更新 `README.md` / `docs/layouts.md`：使用 settlement statement 的示例与配置字段。

## 迁移说明
- 这是破坏性变更：paymentstatement API 被移除。
- 用户若仍需类似功能，改用 settlement statement。

## 测试与验收
- `go test ./...` 通过。
- `cmd/generate-samples` 生成 settlement statement 样例。
- invoice 输出不受影响。
- settlement statement 文案与 invoice 有明确区分（如 Bill To -> Recipient、Invoice ID -> Statement ID、Invoice Issue Date -> Statement Issue Date、Invoice Period -> Statement Period）。

## TODO
- [x] 删除 paymentstatement 相关代码与入口（core/builder/layouts/README/docs/samples/tests）。
- [x] 新增 settlement statement 对外接口与构造方法。
- [x] 复用 invoice layout，并在文案处切换为 Statement 版本。
- [x] i18n 补齐 Statement 文案（en/ja/zh_cn/zh_tw）。
- [x] 更新 samples 与生成逻辑。
- [x] 更新 README 与 docs/layouts.md。
- [x] 运行测试与样例生成，确认输出无回归。
