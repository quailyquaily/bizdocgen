package builder

import (
	"log"
	"log/slog"

	"github.com/johnfercher/maroto/v2/pkg/components/page"
	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/quailyquaily/bizdocgen/core"
	"github.com/quailyquaily/bizdocgen/i18n"
)

type (
	Config struct {
		FontName       string
		FontNormal     string
		FontItalic     string
		FontBold       string
		FontBoldItalic string

		Lang string

		// Layout names: "classic" (default), "modern", "compact".
		InvoiceLayout          string
		PaymentStatementLayout string
	}

	Builder struct {
		cfg              Config
		i18nBundle       *i18n.I18nBundle
		iParams          *core.InvoiceParams
		psParams         *core.PaymentStatementParams
		Round            int32
		fgColor          *props.Color
		fgSecondaryColor *props.Color
		fgTertiaryColor  *props.Color
		borderColor      *props.Color
	}
)

func NewInvoiceBuilder(cfg Config, params *core.InvoiceParams) (*Builder, error) {
	i18nBundle := i18n.New()
	if cfg.Lang == "" {
		cfg.Lang = "en"
	}
	round := 2
	if params.Currency == "JPY" || params.Currency == "円" {
		round = 0
	}
	return &Builder{
		cfg:              cfg,
		i18nBundle:       i18nBundle,
		iParams:          params,
		Round:            int32(round),
		fgColor:          &props.Color{Red: 50, Green: 50, Blue: 93},
		fgSecondaryColor: &props.Color{Red: 80, Green: 80, Blue: 123},
		fgTertiaryColor:  &props.Color{Red: 120, Green: 120, Blue: 153},
		borderColor:      &props.Color{Red: 210, Green: 210, Blue: 230},
	}, nil
}

func NewInvoiceBuilderFromFile(cfg Config, filename string) (*Builder, error) {
	params := &core.InvoiceParams{}
	if err := params.Load(filename); err != nil {
		return nil, err
	}
	return NewInvoiceBuilder(cfg, params)
}

func NewPaymentStatementBuilder(cfg Config, params *core.PaymentStatementParams) (*Builder, error) {
	i18nBundle := i18n.New()
	if cfg.Lang == "" {
		cfg.Lang = "en"
	}
	round := 2
	if params.Currency == "JPY" || params.Currency == "円" {
		round = 0
	}
	return &Builder{
		cfg:              cfg,
		i18nBundle:       i18nBundle,
		psParams:         params,
		Round:            int32(round),
		fgColor:          &props.Color{Red: 50, Green: 50, Blue: 93},
		fgSecondaryColor: &props.Color{Red: 80, Green: 80, Blue: 123},
		borderColor:      &props.Color{Red: 210, Green: 210, Blue: 230},
	}, nil
}

func NewPaymentStatementBuilderFromFile(cfg Config, filename string) (*Builder, error) {
	params := &core.PaymentStatementParams{}
	if err := params.Load(filename); err != nil {
		return nil, err
	}
	return NewPaymentStatementBuilder(cfg, params)
}

func (b *Builder) GenerateInvoice() ([]byte, error) {
	layout := InvoiceLayoutByName(b.cfg.InvoiceLayout)
	headers, body, err := layout.Build(b)
	if err != nil {
		log.Printf("failed to build invoice layout %q: %v\n", layout.Name(), err)
		return nil, err
	}

	m, err := b.CreateMetricsDecorator(headers)
	if err != nil {
		log.Printf("failed to register header: %v\n", err)
		return nil, err
	}

	footer, err := b.BuildInvoiceFooter()
	if err != nil {
		log.Printf("failed to build invoice footer: %v\n", err)
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

func (b *Builder) GeneratePaymentStatement() ([]byte, error) {
	layout := PaymentStatementLayoutByName(b.cfg.PaymentStatementLayout)
	headers, body, err := layout.Build(b)
	if err != nil {
		log.Printf("failed to build payment statement layout %q: %v\n", layout.Name(), err)
		return nil, err
	}

	m, err := b.CreateMetricsDecorator(headers)
	if err != nil {
		log.Printf("failed to register header: %v\n", err)
		return nil, err
	}

	newPage := page.New()
	newPage.Add(body...)

	m.AddPages(newPage)

	return b.getBytesFromMaroto(m)
}

func (b *Builder) getBytesFromMaroto(maroto marotoCore.Maroto) ([]byte, error) {
	document, err := maroto.Generate()
	if err != nil {
		slog.Error("failed to generate document from maroto", "error", err)
		return nil, err
	}

	bytes := document.GetBytes()
	return bytes, nil
}
