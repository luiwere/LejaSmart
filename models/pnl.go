package models

type PnLSummary struct {
    TotalIncome   float64 `json:"total_income"`
    TotalExpenses float64 `json:"total_expenses"`
    Profit        float64 `json:"profit"`
    From          string  `json:"from"`
    To            string  `json:"to"`
}
