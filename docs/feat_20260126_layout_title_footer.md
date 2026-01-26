# 需求整理：发票模板标题与页脚

## 背景
为所有 invoice 模板统一补充首页标题区内容，并在页脚显示页码与发票编号，保持布局一致性与可定制性。

## 需求 1：首页头部标题区
- **适用范围**：所有 invoice 模板（第一页）。
- **默认文案**：
  - 大标题：`Invoice`
  - 小字说明：`This is an invoice hint`
- **文案来源**：
  - 大标题：`InvoiceParams.Doc.Title`（默认值为 `Invoice`）。
  - 小字说明：`InvoiceParams.Doc.Description`（默认值为 `This is an invoice hint`）。
- **版式要求**：所有模板首页头部以相同格式呈现上述两行文字。

## 需求 2：页脚信息
- **适用范围**：所有 invoice 模板（每一页）。
- **布局**：
  - 左侧：显示 Invoice ID（来源：`InvoiceParams.ID`）。
  - 右侧：显示页码，格式为 `当前页 / 总页`（示例：`1 / 2`）。

## 统一排版规范（拟定）
- **字体**：沿用当前模板默认字体（不新增字体族），确保所有 invoice 模板一致。
- **标题区**：
  - 大标题：Bold，`Size=18`，左对齐。
  - 小字说明：Regular，`Size=9`，左对齐。
  - 标题与说明的垂直间距：`1.25pt`。
  - 标题区与正文的垂直间距：`0.5pt`。
- **页脚**：
  - 字号：Regular，`Size=9`。
  - 左右对齐：左侧 Invoice ID，右侧页码。
