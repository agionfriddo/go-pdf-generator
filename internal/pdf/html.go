package pdf

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/url"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/agionfriddo/pdf-poc/internal/models"
)

//go:embed templates/*.tmpl templates/partials/*.tmpl
var paystubTemplateFS embed.FS

type htmlRenderOptions struct {
	TitleSuffix string
	Options     RenderOptions
	Payroll     models.Payroll
}

func renderHTML(p models.Payroll, opts RenderOptions, titleSuffix string) ([]byte, error) {

	tpl, err := template.ParseFS(paystubTemplateFS, "templates/*.tmpl", "templates/partials/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, htmlRenderOptions{TitleSuffix: titleSuffix, Options: opts, Payroll: p}); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}
	return buf.Bytes(), nil
}

func htmlToPDF(htmlBytes []byte) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var pdfData []byte
	dataURL := "data:text/html," + url.PathEscape(string(htmlBytes))

	err := chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.5).
				WithPaperHeight(11).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfData = buf
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("chromedp print to pdf: %w (is Chrome/Chromium installed?)", err)
	}
	return pdfData, nil
}

// GeneratePayrollPDFViaHTML renders the default pay stub using HTML->PDF pipeline.
func GeneratePayrollPDFViaHTML(p models.Payroll) ([]byte, error) {
	html, err := renderHTML(p, RenderOptions{ShowEmployerEmail: false, ShowAddresses: false, ShowBenefits: false}, "")
	if err != nil {
		return nil, err
	}
	return htmlToPDF(html)
}

// GenerateCanadaPayrollPDFViaHTML renders the CA pay stub using HTML->PDF pipeline.
func GenerateCanadaPayrollPDFViaHTML(p models.Payroll) ([]byte, error) {
	html, err := renderHTML(p, RenderOptions{ShowEmployerEmail: true, ShowAddresses: true, ShowBenefits: true}, "Canada")
	if err != nil {
		return nil, err
	}
	return htmlToPDF(html)
}
