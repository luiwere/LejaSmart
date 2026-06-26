package models

type PnLSummary struct {
	TotalRevenue  float64 `json:"total_revenue"`
	TotalCOGS     float64 `json:"total_cogs"`
	TotalExpenses float64 `json:"total_expenses"`
	GrossProfit   float64 `json:"gross_profit"`
	NetProfit     float64 `json:"net_profit"`
	From          string  `json:"from"`
	To            string  `json:"to"`
}
