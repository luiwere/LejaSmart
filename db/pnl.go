package db

type PnLSummary struct {
    TotalIncome   float64 `json:"total_income"`
    TotalExpenses float64 `json:"total_expenses"`
    Profit        float64 `json:"profit"`
    From          string  `json:"from"`
    To            string  `json:"to"`
}

func GetPnL(vendorID, from, to string) (PnLSummary, error) {
    var summary PnLSummary
    summary.From = from
    summary.To = to

    // Build date filter
    dateFilter := ""
    args := []interface{}{vendorID}
    if from != "" && to != "" {
        dateFilter = "AND date BETWEEN ? AND ?"
        args = append(args, from, to)
    }

    // Total expenses
    expQuery := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE vendor_id = ? ` + dateFilter
    DB.QueryRow(expQuery, args...).Scan(&summary.TotalExpenses)

    // Total income
    incQuery := `SELECT COALESCE(SUM(amount), 0) FROM income WHERE vendor_id = ? ` + dateFilter
    DB.QueryRow(incQuery, args...).Scan(&summary.TotalIncome)

    summary.Profit = summary.TotalIncome - summary.TotalExpenses

    return summary, nil
}
