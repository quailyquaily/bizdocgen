package builder

import (
	"strings"

	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"
)

const (
	LayoutClassic   = "classic"
	LayoutModern    = "modern"
	LayoutCompact   = "compact"
	LayoutSpotlight = "spotlight"
	LayoutLedger    = "ledger"
	LayoutSplit     = "split"
)

type InvoiceLayout interface {
	Name() string
	Build(b *Builder) (header []marotoCore.Row, body []marotoCore.Row, err error)
}

type PaymentStatementLayout interface {
	Name() string
	Build(b *Builder) (header []marotoCore.Row, body []marotoCore.Row, err error)
}

func InvoiceLayoutByName(name string) InvoiceLayout {
	switch normalizeLayoutName(name) {
	case LayoutModern:
		return invoiceLayoutModern{}
	case LayoutCompact:
		return invoiceLayoutCompact{}
	case LayoutSpotlight:
		return invoiceLayoutSpotlight{}
	case LayoutLedger:
		return invoiceLayoutLedger{}
	case LayoutSplit:
		return invoiceLayoutSplit{}
	case LayoutClassic:
		fallthrough
	default:
		return invoiceLayoutClassic{}
	}
}

func PaymentStatementLayoutByName(name string) PaymentStatementLayout {
	switch normalizeLayoutName(name) {
	case LayoutModern:
		return paymentStatementLayoutModern{}
	case LayoutCompact:
		return paymentStatementLayoutCompact{}
	case LayoutSpotlight:
		return paymentStatementLayoutSpotlight{}
	case LayoutLedger:
		return paymentStatementLayoutLedger{}
	case LayoutSplit:
		return paymentStatementLayoutSplit{}
	case LayoutClassic:
		fallthrough
	default:
		return paymentStatementLayoutClassic{}
	}
}

func BuiltinLayoutNames() []string {
	return []string{
		LayoutClassic,
		LayoutModern,
		LayoutCompact,
		LayoutSpotlight,
		LayoutLedger,
		LayoutSplit,
	}
}

func normalizeLayoutName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return LayoutClassic
	}
	return name
}
