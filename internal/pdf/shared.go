package pdf

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"

	"github.com/agionfriddo/pdf-poc/internal/models"
)

// RenderOptions controls conditional sections in the PDF.
type RenderOptions struct {
	TitleSuffix       string
	ShowEmployerEmail bool
	ShowAddresses     bool
	ShowBenefits      bool
}

func newPDF(title string) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 18)
	pdf.Cell(0, 10, title)
	pdf.Ln(12)
	return pdf
}

func renderParties(pdf *gofpdf.Fpdf, p models.Payroll, opts RenderOptions) {
	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(95, 6, fmt.Sprintf("Employer: %s", p.Employer.Name), "0", 0, "L", false, 0, "")
	if opts.ShowEmployerEmail {
		pdf.CellFormat(0, 6, fmt.Sprintf("Email: %s", p.Employer.Email), "0", 1, "L", false, 0, "")
	} else {
		pdf.CellFormat(0, 6, fmt.Sprintf("Employer ID: %s", p.Employer.ID), "0", 1, "L", false, 0, "")
	}
	if opts.ShowAddresses && p.Employer.Address != "" {
		pdf.MultiCell(0, 6, fmt.Sprintf("Employer Address: %s", p.Employer.Address), "0", "L", false)
	}

	pdf.CellFormat(95, 6, fmt.Sprintf("Employee: %s", p.Employee.Name), "0", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Title: %s", p.Employee.Title), "0", 1, "L", false, 0, "")
	if opts.ShowAddresses && p.Employee.Address != "" {
		pdf.MultiCell(0, 6, fmt.Sprintf("Employee Address: %s", p.Employee.Address), "0", "L", false)
	}
	pdf.Ln(4)
}

func renderPayPeriod(pdf *gofpdf.Fpdf, p models.Payroll) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(0, 7, "Pay Period")
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(63, 6, fmt.Sprintf("Start: %s", p.PayPeriod.StartDate), "0", 0, "L", false, 0, "")
	pdf.CellFormat(63, 6, fmt.Sprintf("End: %s", p.PayPeriod.EndDate), "0", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Pay Date: %s", p.PayPeriod.PayDate), "0", 1, "L", false, 0, "")
	pdf.Ln(3)
}

func renderAmountTable(pdf *gofpdf.Fpdf, title string, items []models.LineItem, total float64) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(0, 7, title)
	pdf.Ln(7)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(120, 7, "Description", "1", 0, "L", false, 0, "")
	pdf.CellFormat(0, 7, "Amount", "1", 1, "R", false, 0, "")
	pdf.SetFont("Helvetica", "", 11)
	for _, it := range items {
		pdf.CellFormat(120, 7, it.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(0, 7, fmt.Sprintf("$%.2f", it.Amount), "1", 1, "R", false, 0, "")
	}
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(120, 7, "Total", "1", 0, "L", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("$%.2f", total), "1", 1, "R", false, 0, "")
	pdf.Ln(5)
}

func renderNetPay(pdf *gofpdf.Fpdf, p models.Payroll) {
	pdf.SetFont("Helvetica", "B", 13)
	pdf.CellFormat(120, 8, "Net Pay", "0", 0, "L", false, 0, "")
	pdf.CellFormat(0, 8, fmt.Sprintf("$%.2f", p.NetPay), "0", 1, "R", false, 0, "")
}

func outputBytes(pdf *gofpdf.Fpdf) ([]byte, error) {
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
