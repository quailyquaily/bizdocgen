package builder

const (
	labelSetInvoice   = "invoice"
	labelSetStatement = "statement"
)

func (b *Builder) labelKey(invoiceKey string) string {
	if b.docLabelSet != labelSetStatement {
		return invoiceKey
	}
	switch invoiceKey {
	case "InvoiceID":
		return "StatementID"
	case "InvoiceIssueDate":
		return "StatementIssueDate"
	case "InvoicePeriod":
		return "StatementPeriod"
	case "InvoiceBillTo":
		return "StatementRecipient"
	default:
		return invoiceKey
	}
}
