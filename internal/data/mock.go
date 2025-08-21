package data

import (
	"context"
	"encoding/json"

	"github.com/agionfriddo/pdf-poc/internal/models"
)

// FetchPayroll simulates fetching JSON from a remote API and decoding it.
func FetchPayroll(_ context.Context) (models.Payroll, error) {
	// Simulated JSON payload matching the expected schema
	mockJSON := []byte(`{
		"country": "MX",
		"employer": {"name": "Acme Corp", "email": "hr@acme.test", "id": "EMP-1001", "address": "123 King St W, Toronto, ON"},
		"employee": {"name": "Jane Doe", "title": "Software Engineer", "id": "E-4242", "address": "55 Bloor St E, Toronto, ON"},
		"pay_period": {
			"start_date": "2025-08-01",
			"end_date": "2025-08-15",
			"pay_date": "2025-08-20",
			"pay_period_id": "PP-2025-08-1",
			"pay_period_name": "Aug 1 - Aug 15, 2025",
			"pay_period_type": "biweekly"
		},
		"earnings": {
			"data": [
				{"name": "Regular", "amount": 3200},
				{"name": "Overtime", "amount": 250}
			],
			"total": 3450
		},
		"deductions": {
			"data": [
				{"name": "Federal Tax", "amount": 510},
				{"name": "State Tax", "amount": 180},
				{"name": "401k", "amount": 200}
			],
			"total": 890
		},
		"benefits": {
			"data": [
				{"name": "Health Insurance", "amount": 220},
				{"name": "Dental", "amount": 40}
			],
			"total": 260
		},
		"net_pay": 2560
	}`)

	var payroll models.Payroll
	if err := json.Unmarshal(mockJSON, &payroll); err != nil {
		return models.Payroll{}, err
	}
	return payroll, nil
}
