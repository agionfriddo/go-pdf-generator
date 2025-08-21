package models

// Employer represents employer details.
type Employer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	ID      string `json:"id"`
	Address string `json:"address"`
}

// Employee represents employee details.
type Employee struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	ID      string `json:"id"`
	Address string `json:"address"`
}

// PayPeriod holds pay period metadata.
type PayPeriod struct {
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	PayDate       string `json:"pay_date"`
	PayPeriodID   string `json:"pay_period_id"`
	PayPeriodName string `json:"pay_period_name"`
	PayPeriodType string `json:"pay_period_type"`
}

// LineItem represents a single earning or deduction item.
type LineItem struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

// AmountSection groups line items with a total.
type AmountSection struct {
	Data  []LineItem `json:"data"`
	Total float64    `json:"total"`
}

// Payroll is the root payload representing the provided JSON structure.
type Payroll struct {
	Country    string        `json:"country"`
	Employer   Employer      `json:"employer"`
	Employee   Employee      `json:"employee"`
	PayPeriod  PayPeriod     `json:"pay_period"`
	Earnings   AmountSection `json:"earnings"`
	Deductions AmountSection `json:"deductions"`
	Benefits   AmountSection `json:"benefits"`
	NetPay     float64       `json:"net_pay"`
}
