package builder

import (
	"fmt"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type invoiceLayoutClassic struct{}

func (invoiceLayoutClassic) Name() string { return LayoutClassic }

func (invoiceLayoutClassic) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.iParams == nil {
		return nil, nil, fmt.Errorf("invoice params are nil")
	}

	headers, err := b.BuildInvoiceHeader()
	if err != nil {
		return nil, nil, err
	}

	body := make([]marotoCore.Row, 0, 64)
	body = append(body, b.BuildInvoiceBillTo()...)
	body = append(body, b.BuildInvoiceSummaryRows()...)
	body = append(body, b.BuildInvoiceDetailsRows()...)
	if !b.iParams.Payment.Disabled {
		body = append(body, b.BuildInvoicePaymentRows()...)
	}

	return headers, body, nil
}

type invoiceLayoutModern struct{}

func (invoiceLayoutModern) Name() string { return LayoutModern }

func (invoiceLayoutModern) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.iParams == nil {
		return nil, nil, fmt.Errorf("invoice params are nil")
	}

	headers, err := b.buildInvoiceHeader(4)
	if err != nil {
		return nil, nil, err
	}

	tBillTo := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceBillTo", nil)
	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummary", nil)
	tAmount := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryAmount", nil)
	tVAT := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryVAT", nil)
	tTotal := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryTotalWithTax", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 220, Green: 220, Blue: 220},
	}

	billToCol := col.New(6)
	billToCol.Add(
		text.New(tBillTo, props.Text{Size: 10, Top: 0, Style: fontstyle.Bold, Color: b.fgColor}),
		text.New(b.iParams.BillToCompany, props.Text{Size: 10, Top: 8, Style: fontstyle.Bold, Color: b.fgColor}),
	)
	lines := strings.Split(b.iParams.BillToAddress, "\n")
	for ix, line := range lines {
		billToCol.Add(text.New(line, props.Text{Size: 9, Top: float64(6*ix + 16), Color: b.fgColor}))
	}

	summaryNumbers := b.invoiceSummaryNumbers()
	summaryCol := col.New(6)
	summaryCol.Add(
		text.New(tSummary, props.Text{Size: 10, Top: 0, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		text.New(fmt.Sprintf("%s %s", summaryNumbers.Total, b.iParams.Currency), props.Text{Size: 18, Top: 8, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		text.New(b.iParams.Summary.Title, props.Text{Size: 9, Top: 26, Align: align.Right, Color: b.fgSecondaryColor}),
		text.New(fmt.Sprintf("%s: %s %s", tAmount, summaryNumbers.Subtotal.RoundDown(2), b.iParams.Currency), props.Text{Size: 9, Top: 34, Align: align.Right, Color: b.fgSecondaryColor}),
		text.New(fmt.Sprintf("%s: %s %s", tVAT, summaryNumbers.Tax, b.iParams.Currency), props.Text{Size: 9, Top: 40, Align: align.Right, Color: b.fgSecondaryColor}),
		text.New(fmt.Sprintf("%s: %s %s", tTotal, summaryNumbers.Total, b.iParams.Currency), props.Text{Size: 9, Top: 46, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
	)
	if summaryNumbers.QuoteAmount.IsPositive() && summaryNumbers.QuoteText != "" {
		summaryCol.Add(text.New(summaryNumbers.QuoteText, props.Text{Size: 8, Top: 54, Align: align.Right, Color: b.fgSecondaryColor}))
	}

	rowHeight := float64(58)
	if summaryNumbers.QuoteAmount.IsPositive() && summaryNumbers.QuoteText != "" {
		rowHeight = 62
	}

	body := make([]marotoCore.Row, 0, 64)
	body = append(body,
		row.New(rowHeight).WithStyle(borderBottomStyle).Add(billToCol, summaryCol),
		row.New(6),
	)
	body = append(body, b.BuildInvoiceDetailsRows()...)
	if !b.iParams.Payment.Disabled {
		body = append(body, row.New(6))
		body = append(body, b.BuildInvoicePaymentRows()...)
	}

	return headers, body, nil
}

type invoiceLayoutCompact struct{}

func (invoiceLayoutCompact) Name() string { return LayoutCompact }

func (invoiceLayoutCompact) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.iParams == nil {
		return nil, nil, fmt.Errorf("invoice params are nil")
	}

	headers, err := b.buildInvoiceHeader(2)
	if err != nil {
		return nil, nil, err
	}

	tBillTo := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceBillTo", nil)
	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummary", nil)
	tVAT := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryVAT", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 220, Green: 220, Blue: 220},
	}

	billToCol := col.New(7)
	billToCol.Add(
		text.New(tBillTo, props.Text{Size: 9, Top: 0, Style: fontstyle.Bold, Color: b.fgColor}),
		text.New(b.iParams.BillToCompany, props.Text{Size: 9, Top: 7, Style: fontstyle.Bold, Color: b.fgColor}),
	)
	lines := strings.Split(b.iParams.BillToAddress, "\n")
	for ix, line := range lines {
		billToCol.Add(text.New(line, props.Text{Size: 8, Top: float64(5*ix + 14), Color: b.fgColor}))
	}

	summaryNumbers := b.invoiceSummaryNumbers()
	summaryCol := col.New(5)
	summaryCol.Add(
		text.New(tSummary, props.Text{Size: 9, Top: 0, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		text.New(fmt.Sprintf("%s %s", summaryNumbers.Total, b.iParams.Currency), props.Text{Size: 12, Top: 9, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		text.New(fmt.Sprintf("%s: %s %s", tVAT, summaryNumbers.Tax, b.iParams.Currency), props.Text{Size: 8, Top: 24, Align: align.Right, Color: b.fgSecondaryColor}),
	)
	if summaryNumbers.QuoteAmount.IsPositive() && summaryNumbers.QuoteText != "" {
		summaryCol.Add(text.New(summaryNumbers.QuoteText, props.Text{Size: 7, Top: 30, Align: align.Right, Color: b.fgSecondaryColor}))
	}

	body := make([]marotoCore.Row, 0, 64)
	body = append(body,
		row.New(38).WithStyle(borderBottomStyle).Add(billToCol, summaryCol),
		row.New(4),
	)
	body = append(body, b.BuildInvoiceDetailsRows()...)
	if !b.iParams.Payment.Disabled {
		body = append(body, row.New(4))
		body = append(body, b.BuildInvoicePaymentRows()...)
	}

	return headers, body, nil
}

type invoiceLayoutLedger struct{}

func (invoiceLayoutLedger) Name() string { return LayoutLedger }

func (invoiceLayoutLedger) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.iParams == nil {
		return nil, nil, fmt.Errorf("invoice params are nil")
	}

	headers, err := b.buildInvoiceHeader(4)
	if err != nil {
		return nil, nil, err
	}

	body := make([]marotoCore.Row, 0, 64)
	body = append(body, b.BuildInvoiceBillTo()...)
	body = append(body, row.New(4))
	body = append(body, b.BuildInvoiceDetailsRows()...)
	body = append(body, row.New(6))
	body = append(body, b.BuildInvoiceSummaryRows()...)
	if !b.iParams.Payment.Disabled {
		body = append(body, row.New(4))
		body = append(body, b.BuildInvoicePaymentRows()...)
	}

	return headers, body, nil
}

type invoiceLayoutSpotlight struct{}

func (invoiceLayoutSpotlight) Name() string { return LayoutSpotlight }

func (invoiceLayoutSpotlight) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.iParams == nil {
		return nil, nil, fmt.Errorf("invoice params are nil")
	}

	headers, err := b.buildInvoiceHeader(2)
	if err != nil {
		return nil, nil, err
	}

	tAmount := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryAmount", nil)
	tVAT := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryVAT", nil)
	tTotal := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryTotalWithTax", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 220, Green: 220, Blue: 220},
	}

	summaryNumbers := b.invoiceSummaryNumbers()

	spotlightCol := col.New(12)
	spotlightCol.Add(
		text.New(tTotal, props.Text{Size: 10, Top: 2, Align: align.Center, Color: b.fgSecondaryColor}),
		text.New(fmt.Sprintf("%s %s", summaryNumbers.Total, b.iParams.Currency), props.Text{Size: 22, Top: 10, Align: align.Center, Style: fontstyle.Bold, Color: b.fgColor}),
	)
	if b.iParams.Summary.Title != "" {
		spotlightCol.Add(text.New(b.iParams.Summary.Title, props.Text{Size: 9, Top: 26, Align: align.Center, Color: b.fgSecondaryColor}))
	}

	breakdownRow := row.New(16).WithStyle(borderBottomStyle).Add(
		col.New(4).Add(
			text.New(tAmount, props.Text{Size: 8, Top: 2, Align: align.Center, Color: b.fgSecondaryColor}),
			text.New(fmt.Sprintf("%s %s", summaryNumbers.Subtotal.RoundDown(2), b.iParams.Currency), props.Text{Size: 10, Top: 8, Align: align.Center, Style: fontstyle.Bold, Color: b.fgColor}),
		),
		col.New(4).Add(
			text.New(tVAT, props.Text{Size: 8, Top: 2, Align: align.Center, Color: b.fgSecondaryColor}),
			text.New(fmt.Sprintf("%s %s", summaryNumbers.Tax, b.iParams.Currency), props.Text{Size: 10, Top: 8, Align: align.Center, Style: fontstyle.Bold, Color: b.fgColor}),
		),
		col.New(4).Add(
			text.New(tTotal, props.Text{Size: 8, Top: 2, Align: align.Center, Color: b.fgSecondaryColor}),
			text.New(fmt.Sprintf("%s %s", summaryNumbers.Total, b.iParams.Currency), props.Text{Size: 10, Top: 8, Align: align.Center, Style: fontstyle.Bold, Color: b.fgColor}),
		),
	)

	body := make([]marotoCore.Row, 0, 64)
	body = append(body,
		row.New(32).WithStyle(borderBottomStyle).Add(spotlightCol),
	)
	if summaryNumbers.QuoteAmount.IsPositive() && summaryNumbers.QuoteText != "" {
		body = append(body, row.New(8).WithStyle(borderBottomStyle).Add(
			text.NewCol(12, summaryNumbers.QuoteText, props.Text{Size: 8, Top: 2, Align: align.Center, Color: b.fgSecondaryColor}),
		))
	}
	body = append(body, breakdownRow, row.New(6))

	body = append(body, b.BuildInvoiceBillTo()...)
	body = append(body, b.BuildInvoiceDetailsRows()...)
	if !b.iParams.Payment.Disabled {
		body = append(body, row.New(6))
		body = append(body, b.BuildInvoicePaymentRows()...)
	}

	return headers, body, nil
}

type invoiceLayoutSplit struct{}

func (invoiceLayoutSplit) Name() string { return LayoutSplit }

func (invoiceLayoutSplit) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.iParams == nil {
		return nil, nil, fmt.Errorf("invoice params are nil")
	}

	headers, err := b.buildInvoiceHeader(2)
	if err != nil {
		return nil, nil, err
	}

	tPayment := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePayment", nil)
	tMethod := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentMethod", nil)
	tPaymentID := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentID", nil)
	tBankName := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankName", nil)
	tBankBranch := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankBranch", nil)
	tBankDepositType := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankDepositType", nil)
	tBankAccount := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankAccount", nil)
	tBankAccountName := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankAccountName", nil)
	tCryptoCurrency := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentCryptoCurrency", nil)
	tCryptoNetwork := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentCryptoNetwork", nil)
	tCryptoAddress := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentCryptoAddress", nil)
	tCryptoMemo := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentCryptoMemo", nil)
	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummary", nil)
	tAmount := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryAmount", nil)
	tVAT := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryVAT", nil)
	tTotal := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryTotalWithTax", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 220, Green: 220, Blue: 220},
	}

	summaryNumbers := b.invoiceSummaryNumbers()

	summaryCol := col.New(6)
	summaryCol.Add(
		text.New(tSummary, props.Text{Size: 10, Top: 0, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		text.New(fmt.Sprintf("%s %s", summaryNumbers.Total, b.iParams.Currency), props.Text{Size: 16, Top: 10, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		text.New(fmt.Sprintf("%s: %s %s", tAmount, summaryNumbers.Subtotal.RoundDown(2), b.iParams.Currency), props.Text{Size: 9, Top: 30, Align: align.Right, Color: b.fgSecondaryColor}),
		text.New(fmt.Sprintf("%s: %s %s", tVAT, summaryNumbers.Tax, b.iParams.Currency), props.Text{Size: 9, Top: 36, Align: align.Right, Color: b.fgSecondaryColor}),
		text.New(fmt.Sprintf("%s: %s %s", tTotal, summaryNumbers.Total, b.iParams.Currency), props.Text{Size: 9, Top: 42, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
	)
	if summaryNumbers.QuoteAmount.IsPositive() && summaryNumbers.QuoteText != "" {
		summaryCol.Add(text.New(summaryNumbers.QuoteText, props.Text{Size: 8, Top: 50, Align: align.Right, Color: b.fgSecondaryColor}))
	}

	body := make([]marotoCore.Row, 0, 64)
	body = append(body, b.BuildInvoiceBillTo()...)
	body = append(body, row.New(4))
	body = append(body, b.BuildInvoiceDetailsRows()...)
	body = append(body, row.New(6))

	if b.iParams.Payment.Disabled {
		body = append(body, row.New(56).WithStyle(borderBottomStyle).Add(col.New(6), summaryCol))
		return headers, body, nil
	}

	paymentCol := col.New(6)
	paymentCol.Add(text.New(tPayment, props.Text{Size: 10, Top: 0, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}))

	lineTop := 12.0
	lineStep := 6.0
	addLine := func(label, value string) {
		if value == "" {
			return
		}
		paymentCol.Add(text.New(fmt.Sprintf("%s: %s", label, value), props.Text{Size: 9, Top: lineTop, Align: align.Left, Color: b.fgColor}))
		lineTop += lineStep
	}

	method := b.iParams.Payment.Method
	if method == "" {
		method = b.defaultInvoicePaymentMethod()
	}
	addLine(tMethod, method)
	addLine(tPaymentID, b.iParams.Payment.PaymentID)
	addLine(tCryptoCurrency, b.iParams.Payment.ReceiveCryptoCurrency)
	addLine(tCryptoNetwork, b.iParams.Payment.ReceiveCryptoNetwork)
	addLine(tCryptoAddress, b.iParams.Payment.ReceiveCryptoAddress)
	addLine(tCryptoMemo, b.iParams.Payment.ReceiveCryptoMemo)
	addLine(tBankName, b.iParams.Payment.ReceiveAccountBank)
	addLine(tBankBranch, b.iParams.Payment.ReceiveAccountBranch)
	addLine(tBankDepositType, b.iParams.Payment.ReceiveDepositType)
	addLine(tBankAccount, b.iParams.Payment.ReceiveAccountNumber)
	addLine(tBankAccountName, b.iParams.Payment.ReceiveAccountName)
	addLine("SWIFT", b.iParams.Payment.ReceiveAccountSwift)
	addLine("Routing Number", b.iParams.Payment.ReceiveAccountRouting)

	body = append(body, row.New(56).WithStyle(borderBottomStyle).Add(paymentCol, summaryCol))

	return headers, body, nil
}
