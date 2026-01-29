package builder

import (
	"log"

	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/quailyquaily/bizdocgen/core"
)

func NewSettlementStatementBuilder(cfg Config, params *core.SettlementStatementParams) (*Builder, error) {
	builder, err := NewInvoiceBuilder(cfg, invoiceParamsFromSettlement(params))
	if err != nil {
		return nil, err
	}
	builder.docLabelSet = labelSetStatement
	return builder, nil
}

func NewSettlementStatementBuilderFromFile(cfg Config, filename string) (*Builder, error) {
	params := &core.SettlementStatementParams{}
	if err := params.Load(filename); err != nil {
		return nil, err
	}
	return NewSettlementStatementBuilder(cfg, params)
}

func (b *Builder) GenerateSettlementStatement() ([]byte, error) {
	layout := InvoiceLayoutByName(b.cfg.SettlementStatementLayout)
	headers, body, err := layout.Build(b)
	if err != nil {
		log.Printf("failed to build settlement statement layout %q: %v\n", layout.Name(), err)
		return nil, err
	}

	m, err := b.CreateMetricsDecorator(headers)
	if err != nil {
		log.Printf("failed to register header: %v\n", err)
		return nil, err
	}

	footer, err := b.BuildInvoiceFooter()
	if err != nil {
		log.Printf("failed to build settlement statement footer: %v\n", err)
		return nil, err
	}
	if err := m.RegisterFooter(footer...); err != nil {
		log.Printf("failed to register footer: %v\n", err)
		return nil, err
	}

	newPage := page.New()
	newPage.Add(body...)

	m.AddPages(newPage)

	return b.getBytesFromMaroto(m)
}

func invoiceParamsFromSettlement(params *core.SettlementStatementParams) *core.InvoiceParams {
	result := core.InvoicePaymentResult(params.Payment.SettlementStatementPaymentResult)
	result.Disabled = !shouldShowSettlementPaymentResult(result)
	return &core.InvoiceParams{
		ID:            params.ID,
		TaxNumber:     params.TaxNumber,
		Date:          params.Date,
		Currency:      params.Currency,
		CompanyName:   params.CompanyName,
		CompanyAddr:   params.CompanyAddr,
		CompanyEmail:  params.CompanyEmail,
		CompanySeal:   params.CompanySeal,
		BillToCompany: params.RecipientCompany,
		BillToAddress: params.RecipientAddress,
		Summary:       params.Summary,
		DetailItems:   params.DetailItems,
		Payment: core.InvoicePayment{
			InvoicePaymentInstruction: core.InvoicePaymentInstruction{Disabled: true},
			InvoicePaymentResult:      result,
		},
		Doc: params.Doc,
	}
}

func shouldShowSettlementPaymentResult(result core.InvoicePaymentResult) bool {
	if result.Disabled {
		return false
	}
	return result.PaymentMethod != "" ||
		!result.Amount.IsZero() ||
		result.Currency != "" ||
		!result.PaidDate.IsZero() ||
		result.TxID != ""
}
