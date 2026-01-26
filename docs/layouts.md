# Layout Styles Plan

## Goals
- Support multiple built-in “layout styles” (template-level changes: section order, grouping, columns).
- Keep the public API backwards compatible: existing code should continue to render the current layout by default.
- Make adding a new layout a small, isolated change (new file + registry entry).

## Proposed API
- Extend `builder.Config` with:
  - `InvoiceLayout string` (`"classic"` default)
  - `PaymentStatementLayout string` (`"classic"` default)
- Provide built-in layout names as exported constants (e.g. `builder.LayoutClassic`, `builder.LayoutModern`, `builder.LayoutCompact`).

## Refactor Approach
1. Add small layout interfaces in `builder/`:
   - `InvoiceLayout` builds `(headerRows, bodyRows)` for invoices.
   - `PaymentStatementLayout` builds `(headerRows, bodyRows)` for payment statements.
2. Add a registry/selector function that maps layout name → implementation and falls back to `"classic"` on unknown names.
3. Update `(*Builder).GenerateInvoice()` / `GeneratePaymentStatement()` to delegate orchestration to the selected layout:
   - get header/body rows from the layout
   - call `CreateMetricsDecorator(header)`
   - add body rows to a page and render

## Built-in Layouts
- `classic`: baseline; matches the original output.
- `modern`: “dashboard” style:
  - invoice: bill-to + summary in a two-column row; details follow; payment section appended.
  - payment statement: payer/payee side-by-side; channel + summary side-by-side; details follow.
- `compact`: denser version of modern; more “one-page friendly”.
- `spotlight`: big total/payment amount upfront, then details (good for quick scanning).
- `ledger`: “accounting” order; details first, totals/summary at the end.
- `split`: details first, then a two-column footer (e.g., payment vs summary).
  - reduced vertical whitespace, smaller section spacing, keeps the same information but more “one-page friendly”.

## Acceptance Criteria
- Default output stays effectively unchanged when `InvoiceLayout` / `PaymentStatementLayout` are empty.
- Selecting `"modern"` or `"compact"` changes grouping/ordering/columns clearly.
- `go test ./...` passes (note: tests currently write `sample-*.pdf` artifacts).

## Usage (Example)
```go
bd, err := builder.NewInvoiceBuilderFromFile(builder.Config{
	Lang:         "en",
	InvoiceLayout: builder.LayoutModern,
}, "./samples/invoice-1.yaml")
```
