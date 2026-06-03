package db

import (
    "github.com/google/uuid"
)

type InventoryItem struct {
    ID        string  `json:"id"`
    VendorID  string  `json:"vendor_id"`
    Name      string  `json:"name"`
    Quantity  float64 `json:"quantity"`
    Unit      string  `json:"unit"`
    UpdatedAt string  `json:"updated_at"`
}

func GetInventory(vendorID string) ([]InventoryItem, error) {
    rows, err := DB.Query(`SELECT id, vendor_id, name, quantity, unit, updated_at FROM inventory WHERE vendor_id = ?`, vendorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []InventoryItem
    for rows.Next() {
        var i InventoryItem
        rows.Scan(&i.ID, &i.VendorID, &i.Name, &i.Quantity, &i.Unit, &i.UpdatedAt)
        items = append(items, i)
    }
    return items, nil
}
