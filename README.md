# bizdocgen

A business document generator created by [Quail](https://quail.ink).

![](https://static.quail.ink/media/qz5uzv5q.webp)

## Usage

### Generate an invoice

```go
package main

import (
	"log"
	"os"

	"github.com/quailyquaily/bizdocgen/builder"
)

func main() {
	bd, err := builder.NewInvoiceBuilderFromFile(builder.Config{
		Lang:          "en",
		InvoiceLayout: builder.LayoutModern, // classic|modern|compact|spotlight|ledger|split
	}, "./samples/invoice-1.yaml")
	if err != nil {
		log.Panic(err)
	}

	buf, err := bd.GenerateInvoice()
	if err != nil {
		log.Panic(err)
	}
	_ = os.WriteFile("invoice.pdf", buf, 0o666)
}
```

### Fonts + language

To render CJK content, configure UTF-8 fonts (the repo includes Noto Sans CJK under `fonts/`):

```go
	bd, _ := builder.NewPaymentStatementBuilderFromFile(
		builder.Config{
			Lang:                   "ja", // also: zh_cn, zh_tw (or zh-CN / zh-TW)
			PaymentStatementLayout: builder.LayoutModern,

		FontName:       "noto-sans-cjk",
		FontNormal:     "./fonts/NotoSansCJK-JP/NotoSansCJKjp-Regular.ttf",
		FontItalic:     "./fonts/NotoSansCJK-JP/NotoSansCJKjp-Italic.ttf",
		FontBold:       "./fonts/NotoSansCJK-JP/NotoSansCJKjp-Bold.ttf",
		FontBoldItalic: "./fonts/NotoSansCJK-JP/NotoSansCJKjp-BoldItalic.ttf",
	},
	"./samples/paymentstatement-1.yaml",
)
```

## Layouts

Select layouts via `builder.Config.InvoiceLayout` / `builder.Config.PaymentStatementLayout`.
Built-ins: `classic`, `modern`, `compact`, `spotlight`, `ledger`, `split`. See `docs/layouts.md`.

## Quote currency reference (optional)

To display an implied exchange rate in invoice summary, set:
- `summary.total_include_tax_quote_amount`
- `summary.total_include_tax_quota_symbol` (alias: `summary.total_include_tax_quote_symbol`)

## Samples

- Inputs (YAML): `samples/`
- Generate local PDFs: `go run ./cmd/generate-samples` → `samples/` (PDFs are ignored by git via `.gitignore`)
