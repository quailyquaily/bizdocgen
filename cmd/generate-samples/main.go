package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/quailyquaily/bizdocgen/builder"
	"github.com/quailyquaily/bizdocgen/core"
)

func main() {
	repoRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(repoRoot, "samples"), 0o777); err != nil {
		log.Fatal(err)
	}

	if err := generateInvoices(repoRoot); err != nil {
		log.Fatal(err)
	}
	if err := generatePaymentStatements(repoRoot); err != nil {
		log.Fatal(err)
	}
}

func generateInvoices(repoRoot string) error {
	input := filepath.Join(repoRoot, "sample-params", "invoice-1.yaml")
	params := &core.InvoiceParams{}
	if err := params.Load(input); err != nil {
		return err
	}
	params.CompanySeal = resolveRelativeToInput(input, params.CompanySeal)

	for _, layout := range builder.BuiltinLayoutNames() {
		cfg := builder.Config{
			Lang:          "en",
			InvoiceLayout: layout,
		}
		bd, err := builder.NewInvoiceBuilder(cfg, params)
		if err != nil {
			return err
		}
		buf, err := bd.GenerateInvoice()
		if err != nil {
			return err
		}

		out := filepath.Join(repoRoot, "samples", fmt.Sprintf("invoice-1-%s.pdf", layout))
		if err := os.WriteFile(out, buf, 0o666); err != nil {
			return err
		}
		log.Printf("wrote %s", out)
	}

	if err := generateInvoice2Modern(repoRoot); err != nil {
		return err
	}

	return nil
}

func generateInvoice2Modern(repoRoot string) error {
	input := filepath.Join(repoRoot, "sample-params", "invoice-2.yaml")
	params := &core.InvoiceParams{}
	if err := params.Load(input); err != nil {
		return err
	}
	params.CompanySeal = resolveRelativeToInput(input, params.CompanySeal)

	fontsRoot := filepath.Join(repoRoot, "fonts", "NotoSansCJK-JP")
	cfg := builder.Config{
		Lang:          "ja",
		InvoiceLayout: builder.LayoutModern,

		FontName:       "noto-sans-cjk",
		FontNormal:     filepath.Join(fontsRoot, "NotoSansCJKjp-Regular.ttf"),
		FontItalic:     filepath.Join(fontsRoot, "NotoSansCJKjp-Italic.ttf"),
		FontBold:       filepath.Join(fontsRoot, "NotoSansCJKjp-Bold.ttf"),
		FontBoldItalic: filepath.Join(fontsRoot, "NotoSansCJKjp-BoldItalic.ttf"),
	}

	bd, err := builder.NewInvoiceBuilder(cfg, params)
	if err != nil {
		return err
	}
	buf, err := bd.GenerateInvoice()
	if err != nil {
		return err
	}

	out := filepath.Join(repoRoot, "samples", "invoice-2-modern.pdf")
	if err := os.WriteFile(out, buf, 0o666); err != nil {
		return err
	}
	log.Printf("wrote %s", out)
	return nil
}

func generatePaymentStatements(repoRoot string) error {
	input := filepath.Join(repoRoot, "sample-params", "paymentstatement-1.yaml")
	params := &core.PaymentStatementParams{}
	if err := params.Load(input); err != nil {
		return err
	}
	params.CompanySeal = resolveRelativeToInput(input, params.CompanySeal)

	fontsRoot := filepath.Join(repoRoot, "fonts", "NotoSansCJK-JP")

	for _, layout := range builder.BuiltinLayoutNames() {
		cfg := builder.Config{
			Lang:                   "ja",
			PaymentStatementLayout: layout,

			FontName:       "noto-sans-cjk",
			FontNormal:     filepath.Join(fontsRoot, "NotoSansCJKjp-Regular.ttf"),
			FontItalic:     filepath.Join(fontsRoot, "NotoSansCJKjp-Italic.ttf"),
			FontBold:       filepath.Join(fontsRoot, "NotoSansCJKjp-Bold.ttf"),
			FontBoldItalic: filepath.Join(fontsRoot, "NotoSansCJKjp-BoldItalic.ttf"),
		}

		bd, err := builder.NewPaymentStatementBuilder(cfg, params)
		if err != nil {
			return err
		}
		buf, err := bd.GeneratePaymentStatement()
		if err != nil {
			return err
		}

		out := filepath.Join(repoRoot, "samples", fmt.Sprintf("paymentstatement-1-%s.pdf", layout))
		if err := os.WriteFile(out, buf, 0o666); err != nil {
			return err
		}
		log.Printf("wrote %s", out)
	}

	return nil
}

func resolveRelativeToInput(inputPath, maybeRelativePath string) string {
	if maybeRelativePath == "" {
		return ""
	}
	if filepath.IsAbs(maybeRelativePath) {
		return maybeRelativePath
	}
	return filepath.Clean(filepath.Join(filepath.Dir(inputPath), maybeRelativePath))
}
