package builder

import (
	"fmt"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type paymentStatementLayoutClassic struct{}

func (paymentStatementLayoutClassic) Name() string { return LayoutClassic }

func (paymentStatementLayoutClassic) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.psParams == nil {
		return nil, nil, fmt.Errorf("payment statement params are nil")
	}

	headers, err := b.BuildPsHeader()
	if err != nil {
		return nil, nil, err
	}

	body := make([]marotoCore.Row, 0, 64)
	body = append(body, b.BuildPsPayer()...)
	body = append(body, b.BuildPsPayee()...)
	body = append(body, b.BuildPsChannelRows()...)
	body = append(body, b.BuildPsSummaryRows()...)
	body = append(body, b.BuildPsDetailsRows()...)

	return headers, body, nil
}

type paymentStatementLayoutModern struct{}

func (paymentStatementLayoutModern) Name() string { return LayoutModern }

func (paymentStatementLayoutModern) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.psParams == nil {
		return nil, nil, fmt.Errorf("payment statement params are nil")
	}

	headers, err := b.BuildPsHeader()
	if err != nil {
		return nil, nil, err
	}

	tPayer := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayer", nil)
	tPayee := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayee", nil)
	tName := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserName", nil)
	tAddress := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserAddress", nil)
	tTaxID := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserTaxID", nil)
	tContact := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserContact", nil)

	tChannelTitle := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannelTitle", nil)
	tChannel := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannel", nil)
	tTxID := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannelTxID", nil)

	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummary", nil)
	tRevenue := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryRevenue", nil)
	tWithholdingTax := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementWithholdingTax", nil)
	tNetAmount := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryNetAmount", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 220, Green: 220, Blue: 220},
	}

	body := make([]marotoCore.Row, 0, 64)

	body = append(body,
		row.New(14).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, tPayer, props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(6, tPayee, props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(7).Add(
			text.NewCol(2, tName, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.Name, props.Text{Size: 9, Top: 2, Align: align.Right}),
			text.NewCol(2, tName, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.Name, props.Text{Size: 9, Top: 2, Align: align.Right}),
		),
		row.New(7).Add(
			text.NewCol(2, tAddress, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.Address, props.Text{Size: 9, Top: 2, Align: align.Right}),
			text.NewCol(2, tAddress, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.Address, props.Text{Size: 9, Top: 2, Align: align.Right}),
		),
		row.New(7).Add(
			text.NewCol(2, tTaxID, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.TaxNumber, props.Text{Size: 9, Top: 2, Align: align.Right}),
			text.NewCol(2, tTaxID, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.TaxNumber, props.Text{Size: 9, Top: 2, Align: align.Right}),
		),
		row.New(7).WithStyle(borderBottomStyle).Add(
			text.NewCol(2, tContact, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.Contact, props.Text{Size: 9, Top: 2, Align: align.Right}),
			text.NewCol(2, tContact, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.Contact, props.Text{Size: 9, Top: 2, Align: align.Right}),
		),
		row.New(8),
	)

	summaryNumbers := b.paymentStatementSummaryNumbers()
	channelCol := col.New(6)
	channelCol.Add(
		text.New(tChannelTitle, props.Text{Size: 11, Top: 0, Style: fontstyle.Bold}),
		text.New(fmt.Sprintf("%s: %s", tChannel, b.psParams.PaymentChannel), props.Text{Size: 9, Top: 10}),
		text.New(fmt.Sprintf("%s: %s", tTxID, b.psParams.PaymentTxID), props.Text{Size: 9, Top: 16}),
	)

	summaryCol := col.New(6)
	summaryCol.Add(
		text.New(tSummary, props.Text{Size: 11, Top: 0, Align: align.Right, Style: fontstyle.Bold}),
		text.New(fmt.Sprintf("%s %s", summaryNumbers.NetAmount.Round(b.Round), b.psParams.Currency), props.Text{Size: 16, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		text.New(fmt.Sprintf("%s: %s %s", tRevenue, summaryNumbers.Revenue.Round(b.Round), b.psParams.Currency), props.Text{Size: 9, Top: 28, Align: align.Right}),
		text.New(fmt.Sprintf("%s: -%s %s", tWithholdingTax, summaryNumbers.WithholdingTax.Round(b.Round), b.psParams.Currency), props.Text{Size: 9, Top: 34, Align: align.Right}),
		text.New(fmt.Sprintf("%s: %s %s", tNetAmount, summaryNumbers.NetAmount.Round(b.Round), b.psParams.Currency), props.Text{Size: 9, Top: 40, Align: align.Right, Style: fontstyle.Bold}),
	)

	body = append(body,
		row.New(46).WithStyle(borderBottomStyle).Add(channelCol, summaryCol),
		row.New(8),
	)
	body = append(body, b.BuildPsDetailsRows()...)

	return headers, body, nil
}

type paymentStatementLayoutCompact struct{}

func (paymentStatementLayoutCompact) Name() string { return LayoutCompact }

func (paymentStatementLayoutCompact) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.psParams == nil {
		return nil, nil, fmt.Errorf("payment statement params are nil")
	}

	headers, err := b.BuildPsHeader()
	if err != nil {
		return nil, nil, err
	}

	tPayer := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayer", nil)
	tPayee := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayee", nil)
	tName := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserName", nil)
	tTaxID := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserTaxID", nil)

	tChannelTitle := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannelTitle", nil)
	tChannel := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannel", nil)

	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummary", nil)
	tNetAmount := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryNetAmount", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 220, Green: 220, Blue: 220},
	}

	body := make([]marotoCore.Row, 0, 64)
	body = append(body,
		row.New(12).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, tPayer, props.Text{Size: 11, Top: 7, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(6, tPayee, props.Text{Size: 11, Top: 7, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(6).Add(
			text.NewCol(2, tName, props.Text{Size: 8, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.Name, props.Text{Size: 8, Top: 2, Align: align.Right}),
			text.NewCol(2, tName, props.Text{Size: 8, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.Name, props.Text{Size: 8, Top: 2, Align: align.Right}),
		),
		row.New(6).WithStyle(borderBottomStyle).Add(
			text.NewCol(2, tTaxID, props.Text{Size: 8, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.TaxNumber, props.Text{Size: 8, Top: 2, Align: align.Right}),
			text.NewCol(2, tTaxID, props.Text{Size: 8, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.TaxNumber, props.Text{Size: 8, Top: 2, Align: align.Right}),
		),
		row.New(6),
	)

	summaryNumbers := b.paymentStatementSummaryNumbers()
	body = append(body,
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, tChannelTitle, props.Text{Size: 10, Top: 9, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(6, tSummary, props.Text{Size: 10, Top: 9, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(8).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, fmt.Sprintf("%s: %s", tChannel, b.psParams.PaymentChannel), props.Text{Size: 8, Top: 2, Align: align.Left}),
			text.NewCol(6, fmt.Sprintf("%s: %s %s", tNetAmount, summaryNumbers.NetAmount.Round(b.Round), b.psParams.Currency), props.Text{Size: 8, Top: 2, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(6),
	)
	body = append(body, b.BuildPsDetailsRows()...)

	return headers, body, nil
}

type paymentStatementLayoutLedger struct{}

func (paymentStatementLayoutLedger) Name() string { return LayoutLedger }

func (paymentStatementLayoutLedger) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.psParams == nil {
		return nil, nil, fmt.Errorf("payment statement params are nil")
	}

	headers, err := b.BuildPsHeader()
	if err != nil {
		return nil, nil, err
	}

	body := make([]marotoCore.Row, 0, 64)
	body = append(body, b.BuildPsPayer()...)
	body = append(body, b.BuildPsPayee()...)
	body = append(body, b.BuildPsChannelRows()...)
	body = append(body, row.New(6))
	body = append(body, b.BuildPsDetailsRows()...)
	body = append(body, row.New(8))
	body = append(body, b.BuildPsSummaryRows()...)

	return headers, body, nil
}

type paymentStatementLayoutSpotlight struct{}

func (paymentStatementLayoutSpotlight) Name() string { return LayoutSpotlight }

func (paymentStatementLayoutSpotlight) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.psParams == nil {
		return nil, nil, fmt.Errorf("payment statement params are nil")
	}

	headers, err := b.BuildPsHeader()
	if err != nil {
		return nil, nil, err
	}

	tRevenue := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryRevenue", nil)
	tWithholdingTax := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementWithholdingTax", nil)
	tNetAmount := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryNetAmount", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 220, Green: 220, Blue: 220},
	}

	summaryNumbers := b.paymentStatementSummaryNumbers()
	spotlightCol := col.New(12)
	spotlightCol.Add(
		text.New(tNetAmount, props.Text{Size: 10, Top: 2, Align: align.Center}),
		text.New(fmt.Sprintf("%s %s", summaryNumbers.NetAmount.Round(b.Round), b.psParams.Currency), props.Text{Size: 20, Top: 10, Align: align.Center, Style: fontstyle.Bold}),
		text.New(fmt.Sprintf("%s: %s %s", tRevenue, summaryNumbers.Revenue.Round(b.Round), b.psParams.Currency), props.Text{Size: 9, Top: 28, Align: align.Center}),
		text.New(fmt.Sprintf("%s: -%s %s", tWithholdingTax, summaryNumbers.WithholdingTax.Round(b.Round), b.psParams.Currency), props.Text{Size: 9, Top: 34, Align: align.Center}),
	)

	tPayer := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayer", nil)
	tPayee := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayee", nil)
	tName := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserName", nil)
	tTaxID := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserTaxID", nil)

	body := make([]marotoCore.Row, 0, 64)
	body = append(body,
		row.New(44).WithStyle(borderBottomStyle).Add(spotlightCol),
		row.New(8),
	)
	body = append(body, b.BuildPsChannelRows()...)
	body = append(body, row.New(6))
	body = append(body, b.BuildPsDetailsRows()...)
	body = append(body,
		row.New(10),
		row.New(12).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, tPayer, props.Text{Size: 11, Top: 7, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(6, tPayee, props.Text{Size: 11, Top: 7, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(7).Add(
			text.NewCol(2, tName, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.Name, props.Text{Size: 9, Top: 2, Align: align.Right}),
			text.NewCol(2, tName, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.Name, props.Text{Size: 9, Top: 2, Align: align.Right}),
		),
		row.New(7).WithStyle(borderBottomStyle).Add(
			text.NewCol(2, tTaxID, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.TaxNumber, props.Text{Size: 9, Top: 2, Align: align.Right}),
			text.NewCol(2, tTaxID, props.Text{Size: 9, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.TaxNumber, props.Text{Size: 9, Top: 2, Align: align.Right}),
		),
	)

	return headers, body, nil
}

type paymentStatementLayoutSplit struct{}

func (paymentStatementLayoutSplit) Name() string { return LayoutSplit }

func (paymentStatementLayoutSplit) Build(b *Builder) ([]marotoCore.Row, []marotoCore.Row, error) {
	if b.psParams == nil {
		return nil, nil, fmt.Errorf("payment statement params are nil")
	}

	headers, err := b.BuildPsHeader()
	if err != nil {
		return nil, nil, err
	}

	tPayer := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayer", nil)
	tPayee := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayee", nil)
	tName := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserName", nil)
	tTaxID := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserTaxID", nil)

	tChannelTitle := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannelTitle", nil)
	tChannel := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannel", nil)
	tTxID := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannelTxID", nil)

	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummary", nil)
	tRevenue := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryRevenue", nil)
	tWithholdingTax := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementWithholdingTax", nil)
	tNetAmount := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryNetAmount", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 220, Green: 220, Blue: 220},
	}

	body := make([]marotoCore.Row, 0, 64)
	body = append(body,
		row.New(12).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, tPayer, props.Text{Size: 11, Top: 7, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(6, tPayee, props.Text{Size: 11, Top: 7, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(6).Add(
			text.NewCol(2, tName, props.Text{Size: 8, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.Name, props.Text{Size: 8, Top: 2, Align: align.Right}),
			text.NewCol(2, tName, props.Text{Size: 8, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.Name, props.Text{Size: 8, Top: 2, Align: align.Right}),
		),
		row.New(6).WithStyle(borderBottomStyle).Add(
			text.NewCol(2, tTaxID, props.Text{Size: 8, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payer.TaxNumber, props.Text{Size: 8, Top: 2, Align: align.Right}),
			text.NewCol(2, tTaxID, props.Text{Size: 8, Top: 2, Align: align.Left}),
			text.NewCol(4, b.psParams.Payee.TaxNumber, props.Text{Size: 8, Top: 2, Align: align.Right}),
		),
		row.New(6),
	)
	body = append(body, b.BuildPsDetailsRows()...)
	body = append(body, row.New(6))

	summaryNumbers := b.paymentStatementSummaryNumbers()
	channelCol := col.New(6)
	channelCol.Add(
		text.New(tChannelTitle, props.Text{Size: 10, Top: 0, Style: fontstyle.Bold}),
		text.New(fmt.Sprintf("%s: %s", tChannel, b.psParams.PaymentChannel), props.Text{Size: 9, Top: 12}),
	)
	if b.psParams.PaymentTxID != "" {
		channelCol.Add(text.New(fmt.Sprintf("%s: %s", tTxID, b.psParams.PaymentTxID), props.Text{Size: 9, Top: 18}))
	}

	summaryCol := col.New(6)
	summaryCol.Add(
		text.New(tSummary, props.Text{Size: 10, Top: 0, Align: align.Right, Style: fontstyle.Bold}),
		text.New(fmt.Sprintf("%s %s", summaryNumbers.NetAmount.Round(b.Round), b.psParams.Currency), props.Text{Size: 14, Top: 10, Align: align.Right, Style: fontstyle.Bold}),
		text.New(fmt.Sprintf("%s: %s %s", tRevenue, summaryNumbers.Revenue.Round(b.Round), b.psParams.Currency), props.Text{Size: 9, Top: 28, Align: align.Right}),
		text.New(fmt.Sprintf("%s: -%s %s", tWithholdingTax, summaryNumbers.WithholdingTax.Round(b.Round), b.psParams.Currency), props.Text{Size: 9, Top: 34, Align: align.Right}),
		text.New(fmt.Sprintf("%s: %s %s", tNetAmount, summaryNumbers.NetAmount.Round(b.Round), b.psParams.Currency), props.Text{Size: 9, Top: 40, Align: align.Right, Style: fontstyle.Bold}),
	)

	body = append(body, row.New(46).WithStyle(borderBottomStyle).Add(channelCol, summaryCol))

	return headers, body, nil
}
