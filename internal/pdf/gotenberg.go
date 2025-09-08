package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/agionfriddo/pdf-poc/internal/models"
)

const defaultGotenbergTimeout = 20 * time.Second

func convertHTMLWithGotenberg(ctx context.Context, baseURL string, html []byte) ([]byte, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Attach index.html
	fw, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := fw.Write(html); err != nil {
		return nil, fmt.Errorf("write html: %w", err)
	}

	// Optional print options
	_ = writer.WriteField("paperWidth", "8.5")
	_ = writer.WriteField("paperHeight", "11")
	_ = writer.WriteField("printBackground", "true")
	_ = writer.WriteField("marginTop", "0.5")
	_ = writer.WriteField("marginBottom", "0.5")
	_ = writer.WriteField("marginLeft", "0.5")
	_ = writer.WriteField("marginRight", "0.5")

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/forms/chromium/convert/html", &body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gotenberg request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gotenberg error: status=%d body=%s", resp.StatusCode, string(b))
	}
	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return pdfBytes, nil
}

// GeneratePayrollPDFViaGotenberg converts the default pay stub using a Gotenberg server.
func GeneratePayrollPDFViaGotenberg(p models.Payroll, baseURL string) ([]byte, error) {
	html, err := renderHTML(p, RenderOptions{ShowEmployerEmail: false, ShowAddresses: false, ShowBenefits: false}, "")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultGotenbergTimeout)
	defer cancel()
	return convertHTMLWithGotenberg(ctx, baseURL, html)
}

// GenerateCanadaPayrollPDFViaGotenberg converts the CA pay stub using a Gotenberg server.
func GenerateCanadaPayrollPDFViaGotenberg(p models.Payroll, baseURL string) ([]byte, error) {
	html, err := renderHTML(p, RenderOptions{ShowEmployerEmail: true, ShowAddresses: true, ShowBenefits: true}, "Canada")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultGotenbergTimeout)
	defer cancel()
	return convertHTMLWithGotenberg(ctx, baseURL, html)
}
