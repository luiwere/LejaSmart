package db

import (
    "github.com/google/uuid"
    "Digiledger/models"
)

func GetInventory(role, vendorID string) ([]models.InventoryItem, error) {
    conn := DBForRole(role)
    rows, err := conn.Query(`SELECT id, vendor_id, name, quantity, unit, updated_at FROM inventory WHERE vendor_id = ?`, vendorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.InventoryItem
    for rows.Next() {
        var i models.InventoryItem
        rows.Scan(&i.ID, &i.VendorID, &i.Name, &i.Quantity, &i.Unit, &i.UpdatedAt)
        items = append(items, i)
    }
    return items, nil
}

func SaveInventoryItem(role, vendorID, name string, quantity float64, unit string) error {
    conn := DBForRole(role)
    id := uuid.New().String()
    _, err := conn.Exec(
        `INSERT OR REPLACE INTO inventory (id, vendor_id, name, quantity, unit, updated_at)
         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
        id, vendorID, name, quantity, unit,
    )
    return err
}