package db

import (
	"Digiledger/models"
)

func GetPnL(role, shopID, vendorID, from, to string) (models.PnLSummary, error) {
	var summary models.PnLSummary
	summary.From = from
	summary.To = to

	conn := DBForRole(role)

	// Build date filter
	dateFilter := ""
	args := []interface{}{shopID}
	if vendorID != "" {
		dateFilter = "AND vendor_id = ?"
		args = append(args, vendorID)
	}
	if from != "" && to != "" {
		dateFilter += " AND date BETWEEN ? AND ?"
		args = append(args, from, to)
	}

	// Total expenses
	expQuery := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE shop_id = ? ` + dateFilter
	conn.QueryRow(expQuery, args...).Scan(&summary.TotalExpenses)

	// Total revenue is derived from sales quantity × unit price only
	var salesRevenue float64
	salesQuery := `SELECT COALESCE(SUM(quantity * unit_price), 0), COALESCE(SUM(quantity * COALESCE(unit_cost, 0)), 0) FROM sales WHERE shop_id = ? ` + dateFilter
	conn.QueryRow(salesQuery, args...).Scan(&salesRevenue, &summary.TotalCOGS)

	summary.TotalRevenue = salesRevenue
	summary.GrossProfit = summary.TotalRevenue - summary.TotalCOGS
	summary.NetProfit = summary.GrossProfit - summary.TotalExpenses

	return summary, nil
}
