package builder

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/shopspring/decimal"
)

const (
	defaultInvoiceDocTitle       = "Invoice"
	defaultInvoiceDocDescription = "This is an invoice hint"
	defaultStatementDocTitle     = "Settlement Statement"
	defaultStatementDocHint      = ""
)

func (b *Builder) BuildInvoiceHeader() ([]marotoCore.Row, error) {
	return b.buildInvoiceHeader(6)
}

func (b *Builder) invoiceDocTitle() string {
	if b.iParams.Doc.Title != "" {
		return b.iParams.Doc.Title
	}
	if b.docLabelSet == labelSetStatement {
		return defaultStatementDocTitle
	}
	return defaultInvoiceDocTitle
}

func (b *Builder) invoiceDocHint() string {
	if b.iParams.Doc.Description != "" {
		return b.iParams.Doc.Description
	}
	if b.docLabelSet == labelSetStatement {
		return defaultStatementDocHint
	}
	return defaultInvoiceDocDescription
}

func (b *Builder) buildInvoiceTitleRows() []marotoCore.Row {
	title := b.invoiceDocTitle()
	hint := b.invoiceDocHint()

	ret := []marotoCore.Row{
		row.New(10).Add(text.NewCol(12, title, props.Text{Size: 20, Top: 0, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor})),
	}

	if hint != "" {
		ret = append(ret, row.New(4).Add(text.NewCol(12, hint, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgTertiaryColor})))
	}
	return ret
}

func (b *Builder) buildInvoiceHeader(spacerHeight float64) ([]marotoCore.Row, error) {
	tInvoiceID := b.i18nBundle.MusT(b.cfg.Lang, b.labelKey("InvoiceID"), nil)
	tTaxID := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceTaxID", nil)
	tIssueDate := b.i18nBundle.MusT(b.cfg.Lang, b.labelKey("InvoiceIssueDate"), nil)
	tPeriod := b.i18nBundle.MusT(b.cfg.Lang, b.labelKey("InvoicePeriod"), nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: b.borderColor,
	}
	leftCol := col.New(6)

	if b.iParams.CompanySeal != "" {
		fd, err := os.Open(b.iParams.CompanySeal)
		if err != nil {
			log.Printf("failed to open seal file: %v\n", err)
			return nil, err
		}
		defer fd.Close()
		buf, err := io.ReadAll(fd)
		if err != nil {
			log.Printf("failed to read seal file: %v\n", err)
			return nil, err
		}

		leftCol.Add(image.NewFromBytes(buf, extension.Png, props.Rect{
			Center:  false,
			Percent: 20,
			Left:    34,
			Top:     7,
		}))
	}

	leftCol.Add(text.New(b.iParams.CompanyName, props.Text{Size: 14, Top: 8, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}))
	lines := strings.Split(b.iParams.CompanyAddr, "\n")
	for ix, line := range lines {
		leftCol.Add(text.New(line, props.Text{Size: 9, Top: float64(6*ix + 16), Align: align.Left, Color: b.fgColor}))
	}
	leftCol.Add(text.New(b.iParams.CompanyEmail, props.Text{Size: 9, Top: float64(6*(len(lines)) + 16), Align: align.Left, Color: b.fgColor}))

	rs := row.New(42).WithStyle(borderBottomStyle).Add(
		leftCol,
		col.New(6).Add(
			text.New(fmt.Sprintf("%s: %s", tInvoiceID, b.iParams.ID), props.Text{Size: 9, Top: 16, Align: align.Right, Color: b.fgColor}),
			text.New(fmt.Sprintf("%s: %s", tTaxID, b.iParams.TaxNumber), props.Text{Size: 9, Top: 22, Align: align.Right, Color: b.fgColor}),
			text.New(fmt.Sprintf("%s: %s", tIssueDate, b.iParams.Date.Format("2006/01/02")), props.Text{Size: 9, Top: 28, Align: align.Right, Color: b.fgColor}),
			text.New(fmt.Sprintf("%s: %s - %s", tPeriod,
				b.iParams.Summary.PeriodStart.Format("2006/01/02"),
				b.iParams.Summary.PeriodEnd.Format("2006/01/02"),
			), props.Text{Size: 9, Top: 34, Align: align.Right, Color: b.fgColor}),
		),
	)

	rows := make([]marotoCore.Row, 0, 4)
	rows = append(rows, b.buildInvoiceTitleRows()...)
	rows = append(rows, rs)
	if spacerHeight > 0 {
		rows = append(rows, row.New(spacerHeight))
	}
	return rows, nil
}

func (b *Builder) BuildInvoiceFooter() ([]marotoCore.Row, error) {
	if b.iParams == nil {
		return nil, fmt.Errorf("invoice params are nil")
	}

	const footerRowHeight = 10.0
	footerRow := row.New(footerRowHeight).Add(
		text.NewCol(6, fmt.Sprintf("%s", b.iParams.ID), props.Text{
			Size:  7,
			Top:   footerRowHeight, // Align with right-bottom page number baseline.
			Align: align.Left,
			Color: b.fgTertiaryColor,
		}),
		col.New(6),
	)

	return []marotoCore.Row{footerRow}, nil
}

func (b *Builder) BuildInvoiceBillTo() []marotoCore.Row {
	tBillTo := b.i18nBundle.MusT(b.cfg.Lang, b.labelKey("InvoiceBillTo"), nil)

	billTo := col.New(8)
	billTo.Add(text.New(b.iParams.BillToCompany, props.Text{Size: 9, Top: float64(0), Style: fontstyle.Bold, Color: b.fgColor}))
	billTo.Add(text.New(b.iParams.BillToAddress, props.Text{Size: 9, Top: float64(6), Color: b.fgColor}))

	return []marotoCore.Row{
		text.NewRow(8, tBillTo, props.Text{Size: 10, Top: 0, Style: fontstyle.Bold, Color: b.fgColor}),
		row.New(12).Add(billTo),
	}
}

func (b *Builder) BuildInvoicePaymentRows() []marotoCore.Row {
	tPayment := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePayment", nil)
	tMethod := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentMethod", nil)
	tBankName := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankName", nil)
	tBankBranch := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankBranch", nil)
	tBankDepositType := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankDepositType", nil)
	tBankAccount := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankAccount", nil)
	tBankAccountName := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankAccountName", nil)
	tCryptoCurrency := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentCryptoCurrency", nil)
	tCryptoNetwork := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentCryptoNetwork", nil)
	tCryptoAddress := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentCryptoAddress", nil)
	tCryptoMemo := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentCryptoMemo", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: b.borderColor,
	}

	rows := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tPayment, props.Text{Size: 10, Top: 8, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}),
			text.NewCol(4, "", props.Text{Size: 10, Top: 8, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		),
	}

	instruction := b.iParams.Payment.InvoicePaymentInstruction
	method := instruction.Method
	if method == "" {
		method = b.defaultInvoicePaymentMethod()
	}
	rows = append(rows, row.New(10).Add(
		col.New(2).Add(
			text.New(tMethod, props.Text{Size: 9, Top: 4, Align: align.Left, Color: b.fgColor}),
		),
		col.New(10).Add(
			text.New(method, props.Text{Size: 9, Top: 4, Align: align.Right, Color: b.fgColor}),
		),
	))

	if instruction.ReceiveCryptoCurrency != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tCryptoCurrency, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveCryptoCurrency, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveCryptoNetwork != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tCryptoNetwork, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveCryptoNetwork, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveCryptoAddress != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tCryptoAddress, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveCryptoAddress, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveCryptoMemo != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tCryptoMemo, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveCryptoMemo, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveAccountBank != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankName, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveAccountBank, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveAccountBranch != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankBranch, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveAccountBranch, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveAccountNumber != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankAccount, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveAccountNumber, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveDepositType != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankDepositType, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveDepositType, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveAccountName != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankAccountName, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveAccountName, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveAccountSwift != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New("SWIFT", props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveAccountSwift, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if instruction.ReceiveAccountRouting != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New("Routing Number", props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(instruction.ReceiveAccountRouting, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}
	return rows
}

func (b *Builder) BuildInvoicePaymentResultRows() []marotoCore.Row {
	tPaymentResult := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentResult", nil)
	tMethod := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentResultMethod", nil)
	tAmount := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentResultAmount", nil)
	tPaidDate := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentResultPaidDate", nil)
	tTxID := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentResultTxID", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: b.borderColor,
	}

	rows := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tPaymentResult, props.Text{Size: 10, Top: 8, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}),
			text.NewCol(4, "", props.Text{Size: 10, Top: 8, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		),
	}

	result := b.iParams.Payment.InvoicePaymentResult
	firstLine := true
	addLine := func(label, value string) {
		if value == "" {
			return
		}
		height := 6.0
		top := 0.0
		if firstLine {
			height = 10
			top = 4
			firstLine = false
		}
		rows = append(rows, row.New(height).Add(
			col.New(2).Add(
				text.New(label, props.Text{Size: 9, Top: top, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(value, props.Text{Size: 9, Top: top, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	addLine(tMethod, result.PaymentMethod)
	if !result.Amount.IsZero() {
		currency := result.Currency
		if currency == "" {
			currency = b.iParams.Currency
		}
		addLine(tAmount, fmt.Sprintf("%s %s", result.Amount.RoundDown(2), currency))
	}
	if !result.PaidDate.IsZero() {
		addLine(tPaidDate, result.PaidDate.Format("2006-01-02"))
	}
	addLine(tTxID, result.TxID)

	return rows
}

func (b *Builder) showInvoicePaymentInstructions() bool {
	if b.iParams == nil {
		return false
	}
	return !b.iParams.Payment.InvoicePaymentInstruction.Disabled
}

func (b *Builder) showInvoicePaymentResult() bool {
	if b.iParams == nil {
		return false
	}
	return !b.iParams.Payment.InvoicePaymentResult.Disabled
}

func (b *Builder) defaultInvoicePaymentMethod() string {
	if b.iParams == nil {
		return ""
	}
	if b.hasInvoiceCryptoPayment() {
		return "Cryptocurrency"
	}
	return "Bank"
}

func (b *Builder) hasInvoiceCryptoPayment() bool {
	if b.iParams == nil {
		return false
	}
	instruction := b.iParams.Payment.InvoicePaymentInstruction
	return instruction.ReceiveCryptoCurrency != "" ||
		instruction.ReceiveCryptoNetwork != "" ||
		instruction.ReceiveCryptoAddress != "" ||
		instruction.ReceiveCryptoMemo != ""
}

func (b *Builder) BuildInvoiceDetailsRows() []marotoCore.Row {
	tDetails := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceDetails", nil)

	colorLink := &props.Color{
		Red:   0,
		Green: 0,
		Blue:  255,
	}

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: b.borderColor,
	}

	rows := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tDetails, props.Text{Size: 10, Top: 8, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}),
			text.NewCol(4, "", props.Text{Size: 10, Top: 8, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		),
	}

	for ix, item := range b.iParams.DetailItems {
		paddingTop := float64(0)
		rowHeight := float64(6)
		if ix == 0 {
			paddingTop = float64(4)
			rowHeight = float64(10)
		}
		r := row.New(rowHeight)
		r.Add(
			col.New(2).Add(
				text.New(item.Date.Format("2006/01/02"), props.Text{Size: 9, Top: paddingTop, Align: align.Left, Color: b.fgColor}),
			),
			col.New(6).Add(
				text.New(item.Title, props.Text{Size: 9, Top: paddingTop, Align: align.Left, Color: b.fgColor}),
			),
		)
		if !item.TotalExcludeTax.IsZero() || !item.TotalIncludeTax.IsZero() {
			if !item.TotalIncludeTax.IsZero() {
				r.Add(
					col.New(4).Add(
						text.New(fmt.Sprintf("%s %s", item.TotalIncludeTax.RoundDown(2), b.iParams.Currency), props.Text{Size: 9, Top: paddingTop, Align: align.Right, Color: b.fgColor}),
					),
				)
			} else {
				r.Add(
					col.New(4).Add(
						text.New(fmt.Sprintf("%s %s", item.TotalExcludeTax.RoundDown(2), b.iParams.Currency), props.Text{Size: 9, Top: paddingTop, Align: align.Right, Color: b.fgColor}),
					),
				)
			}
		}
		rows = append(rows, r)

		if item.Desc != "" {
			r := row.New(6)
			r.Add(
				col.New(2),
			)
			if !item.TotalExcludeTax.IsZero() && !item.Tax.IsZero() {
				r.Add(
					col.New(6).Add(
						text.New(item.Desc, props.Text{Size: 8, Top: 0, Align: align.Left, Color: b.fgSecondaryColor}),
					),
					col.New(4).Add(
						text.New(fmt.Sprintf("VAT: %s %s", item.Tax.RoundDown(2), b.iParams.Currency), props.Text{Size: 8, Top: 0, Align: align.Right, Color: b.fgSecondaryColor}),
					),
				)
			} else {
				r.Add(
					col.New(10).Add(
						text.New(item.Desc, props.Text{Size: 8, Top: 0, Align: align.Left, Color: b.fgSecondaryColor}),
					),
				)
			}
			rows = append(rows, r)
		}
		if item.URL != "" {
			url := item.URL
			rows = append(rows, row.New(6).Add(
				col.New(2),
				col.New(10).Add(
					text.New(item.URL, props.Text{Size: 8, Top: 0, Align: align.Left, Hyperlink: &url, Color: colorLink}),
				),
			))
		} else if len(item.URLs) > 0 {
			for _, url := range item.URLs {
				rows = append(rows, row.New(6).Add(
					col.New(2),
					col.New(10).Add(
						text.New(url, props.Text{Size: 8, Top: 0, Align: align.Left, Hyperlink: &url, Color: colorLink}),
					),
				))
			}
		}
		rows = append(rows, row.New(2))
	}
	return rows
}

func (b *Builder) BuildInvoiceSummaryRows() []marotoCore.Row {
	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummary", nil)
	tAmount := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryAmount", nil)
	tVAT := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryVAT", nil)
	tTotal := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryTotalWithTax", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: b.borderColor,
	}

	summary := b.invoiceSummaryNumbers()

	ret := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tSummary, props.Text{Size: 10, Top: 8, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}),
			text.NewCol(4, tAmount, props.Text{Size: 10, Top: 8, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		),
		row.New(12).Add(
			text.NewCol(8, b.iParams.Summary.Title, props.Text{Size: 9, Top: 4, Align: align.Left, Color: b.fgColor}),
			text.NewCol(4, fmt.Sprintf("%s %s", summary.Subtotal.RoundDown(2), b.iParams.Currency), props.Text{Size: 9, Top: 4, Align: align.Right, Color: b.fgColor}),
		),
		row.New(8).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, tVAT, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			text.NewCol(6, fmt.Sprintf("%s %s", summary.Tax, b.iParams.Currency), props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
		),
		row.New(10).Add(
			text.NewCol(6, tTotal, props.Text{Size: 10, Top: 4, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}),
			text.NewCol(6, fmt.Sprintf("%s %s", summary.Total, b.iParams.Currency), props.Text{Size: 10, Top: 4, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		),
	}
	if summary.QuoteAmount.IsPositive() && summary.QuoteText != "" {
		ret = append(ret, row.New(8).Add(
			text.NewCol(12, summary.QuoteText, props.Text{Size: 8, Top: 2, Align: align.Left, Color: b.fgColor}),
		))
	}
	return ret
}

type invoiceSummaryNumbers struct {
	Subtotal    decimal.Decimal
	Tax         decimal.Decimal
	Total       decimal.Decimal
	QuoteAmount decimal.Decimal
	QuoteSymbol string
	QuoteText   string
}

func (b *Builder) invoiceSummaryNumbers() invoiceSummaryNumbers {
	var total, tax, subtotal, totalJPY decimal.Decimal
	var quoteAmount decimal.Decimal
	var quoteSymbol string
	var quoteText string

	if b.iParams.Summary.TotalExcludeTax.IsPositive() {
		subtotal = b.iParams.Summary.TotalExcludeTax
		if b.iParams.Summary.Tax.IsPositive() {
			tax = b.iParams.Summary.Tax.Round(2)
		} else if b.iParams.Summary.TaxRate.IsPositive() {
			tax = subtotal.Mul(b.iParams.Summary.TaxRate).Round(2)
		}
		total = subtotal.Add(tax).Round(2)
	} else {
		total = b.iParams.Summary.TotalIncludeTax
		subtotal = total.Div(decimal.NewFromFloat(1).Add(b.iParams.Summary.TaxRate)).Round(2)
		tax = total.Sub(subtotal).Round(2)
	}

	quoteSymbol = strings.TrimSpace(b.iParams.Summary.TotalIncludeTaxQuoteSymbol)
	if quoteSymbol == "" {
		quoteSymbol = strings.TrimSpace(b.iParams.Summary.TotalIncludeTaxQuotaSymbol)
	}

	if b.iParams.Summary.TotalIncludeTaxQuoteAmount.IsPositive() && quoteSymbol != "" {
		quoteAmount = b.iParams.Summary.TotalIncludeTaxQuoteAmount
	} else if b.iParams.Summary.TotalIncludeTaxJPY.IsPositive() {
		totalJPY = b.iParams.Summary.TotalIncludeTaxJPY
		quoteAmount = totalJPY
		quoteSymbol = "JPY"
	}

	if quoteAmount.IsPositive() && quoteSymbol != "" && total.IsPositive() {
		quoteAmountRounded := quoteAmount.Round(2)
		quotePerBaseRounded := quoteAmount.Div(total).Round(4)
		if quoteSymbol == "JPY" || quoteSymbol == "円" {
			quoteAmountRounded = quoteAmount.Round(0)
			quotePerBaseRounded = quoteAmount.Div(total).Round(0)
		}

		quoteText = b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryTotalWithTaxQuote", map[string]any{
			"QuoteAmount":  quoteAmountRounded.String(),
			"QuoteSymbol":  quoteSymbol,
			"BaseSymbol":   b.iParams.Currency,
			"QuotePerBase": quotePerBaseRounded.String(),
		})
	}

	return invoiceSummaryNumbers{
		Subtotal:    subtotal,
		Tax:         tax,
		Total:       total,
		QuoteAmount: quoteAmount,
		QuoteSymbol: quoteSymbol,
		QuoteText:   quoteText,
	}
}
