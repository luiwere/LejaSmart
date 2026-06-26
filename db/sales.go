package db

import (
	"Digiledger/models"
	"database/sql"
	"github.com/google/uuid"
)

func GetSales(role, vendorID string) ([]models.Sale, error) {
	conn := DBForRole(role)
	var rows *sql.Rows
	var err error

	if vendorID == "" {
		rows, err = conn.Query(`SELECT id, vendor_id, item_name, quantity, unit_price, unit_cost, date, notes, created_at FROM sales ORDER BY date DESC`)
	} else {
		rows, err = conn.Query(`SELECT id, vendor_id, item_name, quantity, unit_price, unit_cost, date, notes, created_at FROM sales WHERE vendor_id = ? ORDER BY date DESC`, vendorID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var s models.Sale
		rows.Scan(&s.ID, &s.VendorID, &s.ItemName, &s.Quantity, &s.UnitPrice, &s.UnitCost, &s.Date, &s.Notes, &s.CreatedAt)
		sales = append(sales, s)
	}
	return sales, nil
}

func AddSale(role, vendorID, itemName string, quantity, unitPrice, unitCost float64, date, notes string) error {
	conn := DBForRole(role)
	id := uuid.New().String()
	_, err := conn.Exec(
		`INSERT INTO sales (id, vendor_id, item_name, quantity, unit_price, unit_cost, date, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, vendorID, itemName, quantity, unitPrice, unitCost, date, notes,
	)
	return err
}

func DeleteSale(role, id string) error {
	conn := DBForRole(role)
	_, err := conn.Exec(`DELETE FROM sales WHERE id = ?`, id)
	return err
}
