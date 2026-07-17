package db

import (
	"LejaSmart/models"
	"github.com/google/uuid"
)

func GetInventory(role, shopID, vendorID string) ([]models.InventoryItem, error) {
	conn := DBForRole(role)
	rows, err := conn.Query(`SELECT id, vendor_id, name, supplier_name, status, reorder_level, expiry_date, restocked_at, quantity, unit, updated_at FROM inventory WHERE shop_id = ? AND vendor_id = ?`, shopID, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var i models.InventoryItem
		rows.Scan(&i.ID, &i.VendorID, &i.Name, &i.SupplierName, &i.Status, &i.ReorderLevel, &i.ExpiryDate, &i.RestockedAt, &i.Quantity, &i.Unit, &i.UpdatedAt)
		items = append(items, i)
	}
	return items, nil
}

func SaveInventoryItem(role, shopID, vendorID, name, supplierName, status, expiryDate, restockedAt string, reorderLevel, quantity float64, unit string) error {
	conn := DBForRole(role)
	id := uuid.New().String()
	_, err := conn.Exec(
		`INSERT OR REPLACE INTO inventory (id, vendor_id, shop_id, name, supplier_name, status, reorder_level, expiry_date, restocked_at, quantity, unit, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
		id, vendorID, shopID, name, supplierName, status, reorderLevel, expiryDate, restockedAt, quantity, unit,
	)
	return err
}
