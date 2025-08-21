package pdf

import (
	"github.com/agionfriddo/pdf-poc/internal/models"
)

// GenerateCanadaPayrollPDF renders a Canada-specific pay stub PDF using shared helpers.
func GenerateCanadaPayrollPDF(p models.Payroll) ([]byte, error) {
	pdf := newPDF("Pay Statement (Canada)")
	renderParties(pdf, p, RenderOptions{
		ShowEmployerEmail: true,
		ShowAddresses:     true,
		ShowBenefits:      true,
	})
	renderPayPeriod(pdf, p)
	renderAmountTable(pdf, "Earnings", p.Earnings.Data, p.Earnings.Total)
	renderAmountTable(pdf, "Deductions", p.Deductions.Data, p.Deductions.Total)
	// Canada shows benefits
	renderAmountTable(pdf, "Benefits", p.Benefits.Data, p.Benefits.Total)
	renderNetPay(pdf, p)
	return outputBytes(pdf)
}
