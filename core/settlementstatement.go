package core

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type (
	SettlementStatementSummary       = InvoiceSummary
	SettlementStatementDetailItem    = InvoiceDetailItem
	SettlementStatementPaymentResult = InvoicePaymentResult
	SettlementStatementDoc           = InvoiceDoc
)

type SettlementStatementPayment struct {
	SettlementStatementPaymentResult `yaml:"result,omitempty"`
}

type SettlementStatementParams struct {
	ID           string    `yaml:"id"`
	TaxNumber    string    `yaml:"tax_number"`
	Date         time.Time `yaml:"date" time_format:"2006/01/02"`
	Currency     string    `yaml:"currency"`
	CompanyName  string    `yaml:"company_name"`
	CompanyAddr  string    `yaml:"company_address"`
	CompanyEmail string    `yaml:"company_email"`
	CompanySeal  string    `yaml:"company_seal"`

	RecipientCompany string `yaml:"recipient_company"`
	RecipientAddress string `yaml:"recipient_address"`

	// Summary
	Summary SettlementStatementSummary `yaml:"summary"`

	// Details
	DetailItems []SettlementStatementDetailItem `yaml:"detail_items"`

	// Payment Instructions
	Payment SettlementStatementPayment `yaml:"payment"`

	// Doc related info
	Doc SettlementStatementDoc `yaml:"doc"`
}

func (params *SettlementStatementParams) Load(filename string) error {
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
