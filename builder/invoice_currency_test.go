package builder

import (
	"strings"
	"testing"

	"github.com/quailyquaily/bizdocgen/core"
	"github.com/shopspring/decimal"
)

func TestInvoiceSummaryCurrencyOverrideAndFallback(t *testing.T) {
	params := &core.InvoiceParams{
		Currency: "USD",
		Summary: core.InvoiceSummary{
			Currency:                   "EUR",
			Title:                      "Test",
			TotalIncludeTax:            decimal.NewFromInt(100),
			TotalIncludeTaxQuoteAmount: decimal.NewFromInt(13000),
			TotalIncludeTaxQuoteSymbol: "JPY",
		},
	}
	b, err := NewInvoiceBuilder(Config{Lang: "en"}, params)
	if err != nil {
		t.Fatalf("NewInvoiceBuilder: %v", err)
	}

	nums := b.invoiceSummaryNumbers()
	if nums.BaseCurrency != "EUR" {
		t.Fatalf("BaseCurrency = %q, want %q", nums.BaseCurrency, "EUR")
	}
	if !strings.Contains(nums.QuoteText, "1 EUR") {
		t.Fatalf("QuoteText %q does not contain %q", nums.QuoteText, "1 EUR")
	}

	params.Summary.Currency = ""
	nums = b.invoiceSummaryNumbers()
	if nums.BaseCurrency != "USD" {
		t.Fatalf("BaseCurrency = %q, want %q", nums.BaseCurrency, "USD")
	}
	if !strings.Contains(nums.QuoteText, "1 USD") {
		t.Fatalf("QuoteText %q does not contain %q", nums.QuoteText, "1 USD")
	}
}
