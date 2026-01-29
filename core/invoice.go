package core

import (
	"os"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type (
	InvoiceDetailItem struct {
		Date            time.Time       `yaml:"date" time_format:"2006/01/02"`
		Title           string          `yaml:"title"`
		Desc            string          `yaml:"desc"`
		URL             string          `yaml:"url"`
		URLs            []string        `yaml:"urls"`
		Currency        string          `yaml:"currency"`
		TotalExcludeTax decimal.Decimal `yaml:"total_exclude_tax"`
		TotalIncludeTax decimal.Decimal `yaml:"total_include_tax"`
		// TotalIncludeTaxQuoteAmount/TotalIncludeTaxQuoteSymbol provide a reference total in a quote currency,
		// allowing the invoice to display an implied exchange rate for any currency pair.
		TotalIncludeTaxQuoteAmount decimal.Decimal `yaml:"total_include_tax_quote_amount"`
		TotalIncludeTaxQuoteSymbol string          `yaml:"total_include_tax_quote_symbol"`
		Tax                        decimal.Decimal `yaml:"tax"`
	}

	InvoiceSummary struct {
		PeriodStart     time.Time       `yaml:"period_start" time_format:"2006/01/02"`
		PeriodEnd       time.Time       `yaml:"period_end" time_format:"2006/01/02"`
		Title           string          `yaml:"title"`
		Currency        string          `yaml:"currency"`
		TotalExcludeTax decimal.Decimal `yaml:"total_exclude_tax"`
		TotalIncludeTax decimal.Decimal `yaml:"total_include_tax"`
		// TotalIncludeTaxQuoteAmount/TotalIncludeTaxQuoteSymbol provide a reference total in a quote currency,
		// allowing the invoice to display an implied exchange rate for any currency pair.
		TotalIncludeTaxQuoteAmount decimal.Decimal `yaml:"total_include_tax_quote_amount"`
		TotalIncludeTaxQuoteSymbol string          `yaml:"total_include_tax_quote_symbol"`
		// Alias for TotalIncludeTaxQuoteSymbol (kept for backward/typo compatibility).
		TotalIncludeTaxQuotaSymbol string `yaml:"total_include_tax_quota_symbol"`

		// Deprecated: use TotalIncludeTaxQuoteAmount/TotalIncludeTaxQuoteSymbol.
		TotalIncludeTaxJPY decimal.Decimal `yaml:"total_include_tax_jpy"`
		Tax                decimal.Decimal `yaml:"tax"`
		TaxRate            decimal.Decimal `yaml:"tax_rate"`
	}

	InvoicePaymentInstruction struct {
		Disabled bool   `yaml:"disabled"`
		Method   string `yaml:"method"`

		ReceiveAccountBank    string `yaml:"receive_account_bank"`
		ReceiveAccountBranch  string `yaml:"receive_account_branch"`
		ReceiveDepositType    string `yaml:"receive_deposit_type"`
		ReceiveAccountNumber  string `yaml:"receive_account_number"`
		ReceiveAccountName    string `yaml:"receive_account_name"`
		ReceiveAccountRouting string `yaml:"receive_account_routing"`
		ReceiveAccountSwift   string `yaml:"receive_account_swift"`

		ReceiveCryptoCurrency string `yaml:"receive_crypto_currency"`
		ReceiveCryptoNetwork  string `yaml:"receive_crypto_network"`
		ReceiveCryptoAddress  string `yaml:"receive_crypto_address"`
		ReceiveCryptoMemo     string `yaml:"receive_crypto_memo"`
	}

	InvoicePaymentResult struct {
		Disabled      bool            `yaml:"disabled"`
		PaymentMethod string          `yaml:"payment_method"`
		Amount        decimal.Decimal `yaml:"amount_paid"`
		Currency      string          `yaml:"currency"`
		PaidDate      time.Time       `yaml:"paid_date" time_format:"2006-01-02"`
		TxID          string          `yaml:"tx_id"`
	}

	InvoicePayment struct {
		InvoicePaymentInstruction `yaml:"instruction,omitempty"`
		InvoicePaymentResult      `yaml:"result,omitempty"`
	}

	InvoiceDoc struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
	}

	InvoiceParams struct {
		ID           string    `yaml:"id"`
		TaxNumber    string    `yaml:"tax_number"`
		Date         time.Time `yaml:"date" time_format:"2006/01/02"`
		Currency     string    `yaml:"currency"`
		CompanyName  string    `yaml:"company_name"`
		CompanyAddr  string    `yaml:"company_address"`
		CompanyEmail string    `yaml:"company_email"`
		CompanySeal  string    `yaml:"company_seal"`

		BillToCompany string `yaml:"bill_to_company"`
		BillToAddress string `yaml:"bill_to_address"`

		// Summary
		Summary InvoiceSummary `yaml:"summary"`

		// Details
		DetailItems []InvoiceDetailItem `yaml:"detail_items"`

		// Payment Instructions
		Payment InvoicePayment `yaml:"payment"`

		// Doc related info
		Doc InvoiceDoc `yaml:"doc"`
	}
)

func (params *InvoiceParams) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		logrus.WithError(err).Fatalf("failed to read YAML file")
		return err
	}

	if err := yaml.Unmarshal(data, params); err != nil {
		logrus.WithError(err).Fatalf("failed to unmarshal YAML")
		return err
	}

	return nil
}
