package pdf

import (
	"github.com/agionfriddo/pdf-poc/internal/models"
)

// GeneratePayrollPDF renders the default (non-Canada) pay stub PDF using shared helpers.
func GeneratePayrollPDF(p models.Payroll) ([]byte, error) {
	pdf := newPDF("Pay Statement")
	renderParties(pdf, p, RenderOptions{
		ShowEmployerEmail: false,
		ShowAddresses:     false,
		ShowBenefits:      false,
	})
	renderPayPeriod(pdf, p)
	renderAmountTable(pdf, "Earnings", p.Earnings.Data, p.Earnings.Total)
	renderAmountTable(pdf, "Deductions", p.Deductions.Data, p.Deductions.Total)
	renderNetPay(pdf, p)
	return outputBytes(pdf)
}
