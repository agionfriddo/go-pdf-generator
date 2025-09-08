package pdf

import (
	"fmt"

	pw "github.com/playwright-community/playwright-go"

	"github.com/agionfriddo/pdf-poc/internal/models"
)

func htmlToPDFViaPlaywright(htmlBytes []byte) ([]byte, error) {
	if err := pw.Install(); err != nil {
		// Browsers may already be installed in many environments; continue if install fails due to lack of permissions.
		// We will still attempt to run; runtime error will be returned if browsers are unavailable.
		_ = err
	}

	runner, err := pw.Run()
	if err != nil {
		return nil, fmt.Errorf("playwright run: %w", err)
	}
	defer func() { _ = runner.Stop() }()

	browser, err := runner.Chromium.Launch(pw.BrowserTypeLaunchOptions{Headless: pw.Bool(true)})
	if err != nil {
		return nil, fmt.Errorf("launch chromium: %w", err)
	}
	defer func() { _ = browser.Close() }()

	ctx, err := browser.NewContext()
	if err != nil {
		return nil, fmt.Errorf("new context: %w", err)
	}
	defer func() { _ = ctx.Close() }()

	page, err := ctx.NewPage()
	if err != nil {
		return nil, fmt.Errorf("new page: %w", err)
	}

	if err := page.SetContent(string(htmlBytes), pw.PageSetContentOptions{WaitUntil: pw.WaitUntilStateNetworkidle}); err != nil {
		return nil, fmt.Errorf("set content: %w", err)
	}

	pdfBytes, err := page.PDF(pw.PagePdfOptions{
		PrintBackground: pw.Bool(true),
		Format:          pw.String("Letter"),
	})
	if err != nil {
		return nil, fmt.Errorf("page pdf: %w", err)
	}
	return pdfBytes, nil
}

// GeneratePayrollPDFViaPlaywright renders the default pay stub using Playwright.
func GeneratePayrollPDFViaPlaywright(p models.Payroll) ([]byte, error) {
	html, err := renderHTML(p, RenderOptions{ShowEmployerEmail: false, ShowAddresses: false, ShowBenefits: false}, "")
	if err != nil {
		return nil, err
	}
	return htmlToPDFViaPlaywright(html)
}

// GenerateCanadaPayrollPDFViaPlaywright renders the CA pay stub using Playwright.
func GenerateCanadaPayrollPDFViaPlaywright(p models.Payroll) ([]byte, error) {
	html, err := renderHTML(p, RenderOptions{ShowEmployerEmail: true, ShowAddresses: true, ShowBenefits: true}, "Canada")
	if err != nil {
		return nil, err
	}
	return htmlToPDFViaPlaywright(html)
}
